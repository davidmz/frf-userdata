package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/davidmz/frf-userdata/ffapi"
	"github.com/davidmz/mustbe"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

func (a *App) GetPublicParam(r *http.Request) (int, interface{}) {
	pathVars := mux.Vars(r)
	site, username, paramName := pathVars["site"], pathVars["username"], pathVars["paramName"]
	pProps := a.PublicParams[paramName]

	// читаем значение из базы
	vl := pProps.Default
	err := a.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PublicBucketName))
		js := b.Get([]byte(site + "/" + username + "/" + paramName))
		if js != nil {
			return json.Unmarshal(js, &vl)
		}
		return nil
	})

	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, vl
}

func (a *App) PostPublicParam(r *http.Request) (respCode int, respValue interface{}) {
	defer mustbe.Catched(func(err error) {
		a.Log.WithField("error", err).Error(err.Error())
		respCode, respValue = http.StatusInternalServerError, "Internal server error"
	})

	req := &struct {
		AuthToken string          `json:"authToken"`
		Value     json.RawMessage `json:"value"`
	}{}

	pathVars := mux.Vars(r)
	site, username, paramName := pathVars["site"], pathVars["username"], pathVars["paramName"]
	pProps := a.PublicParams[paramName]
	siteInfo := a.Sites[site]

	if err := json.NewDecoder(io.LimitReader(r.Body, int64(pProps.MaxRequestSize))).Decode(req); err != nil {
		return http.StatusBadRequest, err.Error()
	}

	// валидность данных
	val, err := pProps.Schema.Validate(gojsonschema.NewStringLoader(string(req.Value)))
	mustbe.OK(err)

	if !val.Valid() {
		return http.StatusBadRequest, val.Errors()[0].Description()
	}

	api := ffapi.New(siteInfo.APIRoot, req.AuthToken)

	whoami, err := api.WhoAmI()
	mustbe.OK(err)

	if whoami.Users.Username != username {
		// возможно, редактируется группа
		uinfo, err := api.UserInfo(username)
		mustbe.OK(err)
		adminFound := false
		for _, adm := range uinfo.Admins {
			if adm.ID == whoami.Users.ID {
				adminFound = true
				break
			}
		}

		if !adminFound {
			return http.StatusBadRequest, "You can not manage this account"
		}
	}

	// всё в порядке
	mustbe.OK(a.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PublicBucketName))
		return b.Put(
			[]byte(site+"/"+username+"/"+paramName),
			[]byte(req.Value),
		)
	}))

	return http.StatusOK, nil
}

func (a *App) checkVars(h ApiHandler) ApiHandler {
	return func(r *http.Request) (int, interface{}) {
		pathVars := mux.Vars(r)
		site, paramName := pathVars["site"], pathVars["paramName"]

		if _, ok := a.Sites[site]; !ok {
			return http.StatusNotFound, "site not found"
		}

		if _, ok := a.PublicParams[paramName]; !ok {
			return http.StatusNotFound, "invalid parameter name"
		}

		return h(r)
	}
}
