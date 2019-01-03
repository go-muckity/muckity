package muckity

import (
	"os"
	"strconv"
)

type muckityConfig struct {
	mongo mongoConfig
	// TODO: Config the things
}

type MuckityConfig struct {
	done chan interface{}
	muckityConfig
}

// Type implements part of MuckitySystem
func (mc MuckityConfig) Type() string {
	return "systems"
}

// Channels implements part of MuckitySystem
func (mc MuckityConfig) Channels() map[string]chan interface{} {
	slice := make(map[string]chan interface{}, 0)
	slice["done"] = mc.done
	return slice
}

// TODO: use a real config service
func parseConfig(muckCfg *MuckityConfig) {
	if value, ok := os.LookupEnv("MUCKITY_DB_USERNAME"); ok {
		muckCfg.mongo.dbUser = value
	} else {
		muckCfg.mongo.dbUser = "muckity"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_PWD"); ok {
		if muckCfg.mongo.dbUser == "" {
			panic("Empty MUCKITY_DB_USERNAME but MUCKITY_DB_PWD is set!")
		}
		muckCfg.mongo.dbPwd = value
	} else {
		muckCfg.mongo.dbPwd = "muckity"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_HOST"); ok {
		muckCfg.mongo.dbHost = value
	} else {
		muckCfg.mongo.dbHost = "localhost"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_PORT"); ok {
		cfgPort, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			panic(err)
		}
		muckCfg.mongo.dbPort = uint(cfgPort)
	} else {
		muckCfg.mongo.dbPort = 27017
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_NAME"); ok {
		muckCfg.mongo.dbName = value
	} else {
		muckCfg.mongo.dbName = "muckity"
	}
}

func configInit(done chan interface{}) *MuckityConfig {
	mc := new(MuckityConfig)
	parseConfig(mc)
	mc.done = done
	return mc
}
