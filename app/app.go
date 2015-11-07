package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

const (
	PublicBucketName = "PublicVars"
)

type App struct {
	Conf
	DB     *bolt.DB
	Router *mux.Router
	Log    *logrus.Logger
}
