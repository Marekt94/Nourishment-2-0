package utils

import (
	"os"
)

type Error struct {
	Error string `json:"error"`
}

func ReadFile(name string) ([]byte, error) {
	workDir := os.Getenv("WORK_DIR")
	return os.ReadFile(workDir + "/" + name)
}
