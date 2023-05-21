package utils

import (
	"fmt"
	"os"
)

func GetEnv(key, def string, must bool) string {
	res := os.Getenv(key)
	if res == "" {
		if must {
			panic(fmt.Sprintf("env \"%s\" not found", key))
		}
		return def
	}

	return res
}
