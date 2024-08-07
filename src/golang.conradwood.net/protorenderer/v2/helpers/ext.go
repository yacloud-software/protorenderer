package helpers

import (
	"path/filepath"
	"strings"
)

func ChangeExt(filename, new_ext string) string {
	ext := filepath.Ext(filename)
	fname := strings.TrimSuffix(filename, ext)
	res := fname + "." + strings.TrimPrefix(new_ext, ".")
	return res
}
