package app

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/Sirupsen/logrus"
)

type H map[string]interface{}

type ApiHandler func(r *http.Request) (httpCode int, result interface{})

func (a *App) ApiCall(ah ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var (
			httpCode int
			result   interface{}
		)

		func() {
			defer func() {
				// handle panics
				if rec := recover(); rec != nil {
					a.Log.WithField("panic", rec).Error("panic happens")
					if a.Log.Level >= logrus.DebugLevel {
						buf := make([]byte, 1024)
						n := runtime.Stack(buf, false)
						a.Log.WithField("stack", string(buf[:n])).Debug("panic happens")
					}
					httpCode, result = http.StatusInternalServerError, "Internal error"
				}
			}()

			httpCode, result = ah(r)
		}()

		jResult := H{"status": "ok", "data": result}
		if httpCode/100 == http.StatusBadRequest/100 {
			jResult["status"] = "error"
		}
		if httpCode/100 == http.StatusInternalServerError/100 {
			jResult["status"] = "fail"
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		json.NewEncoder(w).Encode(jResult)
	}
}

func (a *App) NotFound(r *http.Request) (int, interface{}) {
	return http.StatusNotFound, "resource not found"
}
