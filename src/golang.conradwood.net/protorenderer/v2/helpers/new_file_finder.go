package helpers

import (
	"golang.conradwood.net/go-easyops/utils"
)

type newFileFinder struct {
	dir      string
	curfiles map[string]bool
}

// find files that are created after some action
func NewFileFinder(dir string) *newFileFinder {
	return &newFileFinder{curfiles: make(map[string]bool), dir: dir}
}

// find new files, relative to dir. internal list of current files will be updated, so that subsequent calls to FindNew() return only files that were created after each call to FindNew(). The first call to FindNew() returns all files.
func (n *newFileFinder) FindNew() ([]string, error) {
	var all_files []string
	err := utils.DirWalk(n.dir, func(root, relfil string) error {
		all_files = append(all_files, relfil)
		return nil
	},
	)
	if err != nil {
		return nil, err
	}
	var new_files []string
	for _, af := range all_files {
		if !n.curfiles[af] {
			new_files = append(new_files, af)
			n.curfiles[af] = true
		}
	}
	return new_files, nil
}
