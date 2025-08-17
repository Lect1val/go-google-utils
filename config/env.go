package config

import (
	"os"
	"strings"
)

var Env string

func init() {
	Env = os.Getenv("ENV")
}

const (
	Local  string = "LOCAL"
	NonPrd string = "NONPRD"
	Prd    string = "PRD"
)

func IsLocalEnv() bool {
	return strings.ToUpper(Env) == Local
}

func IsNonPrdEnv() bool {
	return strings.ToUpper(Env) == NonPrd
}

func IsPrdEnv() bool {
	return strings.ToUpper(Env) == Prd
}
