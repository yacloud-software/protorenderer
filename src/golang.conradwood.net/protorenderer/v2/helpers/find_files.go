package helpers

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/interfaces"
	"strings"
)

func FindProtoFiles(dir string) ([]interfaces.ProtoFile, error) {
	var filenames []string
	extensions := []string{".proto"}
	utils.DirWalk(dir, func(root, relfil string) error {
		for _, ext := range extensions {
			if strings.HasSuffix(relfil, ext) {
				filenames = append(filenames, relfil)
				break
			}
		}
		return nil
	},
	)
	var res []interfaces.ProtoFile
	for _, fn := range filenames {
		fname := dir + "/" + fn
		b, err := utils.ReadFile(fname)
		if err != nil {
			return nil, fmt.Errorf("File %s read error: %s", fname, err)
		}
		spf := &StandardProtoFile{filename: fn, content: b}
		res = append(res, spf)
	}
	return res, nil
}

// returns filenames relative to 'dir'
func FindFiles(dir string, extensions ...string) ([]string, error) {
	var filenames []string
	utils.DirWalk(dir, func(root, relfil string) error {
		for _, ext := range extensions {
			if strings.HasSuffix(relfil, ext) {
				filenames = append(filenames, relfil)
				break
			}
		}
		return nil
	},
	)
	var res []string
	for _, fn := range filenames {
		res = append(res, fn)
	}
	return res, nil
}
