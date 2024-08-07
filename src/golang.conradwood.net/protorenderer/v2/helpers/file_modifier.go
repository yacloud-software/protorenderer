package helpers

import (
	//	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
)

type fileModifier struct {
	filename       string
	headers_to_add []string
	footers_to_add []string
}

func NewFileModifierFromFilename(filename string) *fileModifier {
	return &fileModifier{filename: filename}
}

func (f *fileModifier) AddHeader(header string) {
	f.headers_to_add = append(f.headers_to_add, header)
}
func (f *fileModifier) AddFooter(footer string) {
	f.footers_to_add = append(f.footers_to_add, footer)
}

func (f *fileModifier) Save() error {
	content, err := utils.ReadFile(f.filename)
	if err != nil {
		return err
	}
	content, err = f.modify(content)
	if err != nil {
		return err
	}
	err = utils.WriteFile(f.filename, content)
	if err != nil {
		return err
	}
	//	fmt.Printf("modifier: saved \"%s\"\n", f.filename)
	return nil
}

func (f *fileModifier) modify(content []byte) ([]byte, error) {
	h := strings.Join(f.headers_to_add, "")
	content = append([]byte(h), content...)
	h = strings.Join(f.footers_to_add, "")
	content = append([]byte(h), content...)
	return content, nil
}
