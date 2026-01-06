package utils

import (
	"os"
)

func ReadFile(name string) ([]byte, error) {
	workDir := os.Getenv("WORK_DIR")
	return os.ReadFile(workDir + "/" + name)
}
