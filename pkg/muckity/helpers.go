package muckity

import (
	"github.com/mongodb/mongo-go-driver/mongo/options"
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
