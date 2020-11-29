package util

import (
	"fmt"
	"os"
)

// DebugPrint if DEBUG environment variable is set
func DebugPrint(msg string) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		fmt.Println(msg)
	}
}
