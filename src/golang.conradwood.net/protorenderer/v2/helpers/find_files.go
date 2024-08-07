package helpers

import (
	"golang.conradwood.net/go-easyops/errors"
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
			return nil, errors.Errorf("File %s read error: %s", fname, err)
		}
		spf := &StandardProtoFile{Filename: fn, content: b}
		res = append(res, spf)
	}
	return res, nil
}

// returns filenames relative to 'dir'
func FindFiles(dir string, extensions ...string) ([]string, error) {
	isdir, err := utils.IsDir(dir)
	if err != nil {
		return nil, errors.Errorf("unable to find files: %s", err)
	}
	if !isdir {
		return nil, errors.Errorf("dir \"%s\" does not exist or is not a directory", dir)
	}
	var filenames []string
	utils.DirWalk(dir, func(root, relfil string) error {
		if len(extensions) == 0 {
			filenames = append(filenames, relfil)
			return nil
		}
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
