package app

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/xeipuuv/gojsonschema"
)

type ParamProps struct {
	Default        interface{}
	MaxRequestSize int
	RawSchema      json.RawMessage      `json:"Schema"`
	Schema         *gojsonschema.Schema `json:-`
}

type SiteInfo struct {
	APIRoot string
}

type Conf struct {
	Listen       string
	DBURL        string
	LogLevel     string
	CORSOrigins  []string
	Sites        map[string]*SiteInfo
	PublicParams map[string]*ParamProps
}

func (a *App) LoadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&a.Conf); err != nil {
		return err
	}

	a.DB, err = sql.Open("postgres", a.DBURL)
	if err != nil {
		return err
	}

	if err := a.DB.Ping(); err != nil {
		return err
	}

	a.Router = mux.NewRouter()

	a.Log = logrus.New()
	a.Log.Out = os.Stderr
	a.Log.Formatter = &logrus.TextFormatter{ForceColors: true}
	a.Log.Level = logrus.ErrorLevel
	if ll, err := logrus.ParseLevel(a.LogLevel); err == nil {
		a.Log.Level = ll
	}

	for _, pp := range a.PublicParams {
		loader := gojsonschema.NewStringLoader(string(pp.RawSchema))
		schema, err := gojsonschema.NewSchema(loader)
		if err != nil {
			return err
		}
		pp.Schema = schema
	}

	return nil
}
