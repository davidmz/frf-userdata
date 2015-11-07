package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidmz/frf-userdata/app"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "c", "", "config file")
	flag.Parse()

	if confFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	app := new(app.App)
	if err := app.LoadConfig(confFile); err != nil {
		log.Fatalf("Can not read config file: %v", err)
	}
	defer app.Close()

	app.InitRouter()

	crs := cors.New(cors.Options{
		AllowedOrigins: app.CORSOrigins,
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         86400,
	})

	h := alice.New(
		LoggingHandler(os.Stdout),
		crs.Handler,
	).Then(app.Router)

	s := &http.Server{
		Addr:           app.Listen,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	app.Log.Infof("Starting server at %s", app.Listen)
	app.Log.Fatal(s.ListenAndServe())
}

func LoggingHandler(out io.Writer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(out, h)
	}
}
