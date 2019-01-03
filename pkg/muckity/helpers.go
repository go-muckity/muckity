package muckity

import (
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"os"
	"strconv"
	"strings"
)

func newTrue() *bool {
	b := true
	return &b
}

func newUpsert() options.UpdateOptions {
	o := options.UpdateOptions{Upsert: newTrue()}
	return o
}

func uniqueStrSlice(list []string) []string {
	u := make([]string, 0)
	m := make(map[string]bool)
	for _, val := range list {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// FieldJoin implements an implode function
func FieldJoin(sep string, args ...string) string {
	return strings.Join(args, sep)
}

// TODO: use a real config service
func parseConfig(muckCfg *muckityConfig) {
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
