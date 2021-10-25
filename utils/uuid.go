package utils

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

// return uuid without "-"
func GenRawUuid() string {
	u := uuid.NewV4().String()
	newUuid := strings.Replace(u, "-", "", 4)
	return newUuid
}

func GenUuidFromName(name string) string {
	u := uuid.NewV3(uuid.NamespaceX500, name).String()
	newUuid := strings.Replace(u, "-", "", 4)
	return newUuid
}
