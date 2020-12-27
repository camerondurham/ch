package util

import (
	"os"
	"path/filepath"
)

type Validate interface {
	GetAbs(string) string
	ValidPath(string) bool
}

type Validator struct{}

func (v *Validator) GetAbs(path string) (absPath string) {
	if !filepath.IsAbs(path) {
		absPath, _ = filepath.Abs(path)
	} else {
		absPath = path
	}
	return
}

func (v *Validator) ValidPath(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}
