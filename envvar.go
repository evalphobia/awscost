package awscost

import (
	"os"
)

func getEnvLang() string {
	return os.Getenv("AWSCOST_LANG")
}
