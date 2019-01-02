package main

import (
	"fmt"
	"github.com/tsal/muckity/pkg/muckity"
	"os"
	"strconv"
)

type gameConfig struct {
	dbUser string
	dbPwd string
	dbHost string
	dbPort uint
	dbName string
}

var defaultConfig *gameConfig

func (gc *gameConfig) asURI() string {
	return fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
			gc.dbUser,
			gc.dbPwd,
			gc.dbHost,
			gc.dbPort,
			gc.dbName)
}

func GetConfig() (*gameConfig, error) {
	gc := new(gameConfig)
	if value, ok := os.LookupEnv("MUCKITY_DB_USERNAME"); ok {
		gc.dbUser = value
	} else {
		gc.dbUser = "muckity"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_PWD"); ok {
		if gc.dbUser == "" {
			panic("Empty MUCKITY_DB_USERNAME but MUCKITY_DB_PWD is set!")
		}
		gc.dbPwd = value
	} else {
		gc.dbPwd = "muckity"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_HOST"); ok {
		gc.dbHost = value
	} else {
		gc.dbHost = "localhost"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_PORT"); ok {
		cfgPort, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			panic(err)
		}
		gc.dbPort = uint(cfgPort)
	} else {
		gc.dbPort = 27017
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_NAME"); ok {
		gc.dbName = value
	} else {
		gc.dbName = "muckity"
	}
	return gc, nil
}

func init() {
	var initCfg, err = GetConfig()
	if err != nil {
		panic(err)
	}
	defaultConfig = initCfg
}

func main() {
	uri := defaultConfig.asURI()
	muckity.NewMuckityStorage(uri)
	w, err := muckity.NewWorld("A Brand New World")
	if err != nil {
		panic(err)
	}
	fmt.Println("ID\t: ", w.Name())
}
