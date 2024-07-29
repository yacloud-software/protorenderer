package helpers

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
)

func WriteYaml(filename string, data interface{}) error {
	b, err := utils.MarshalYaml(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for file \"%s\": %w\n", filename, err)
	}
	err = utils.WriteFile(filename, b)
	if err != nil {
		return err
	}
	return nil
}
