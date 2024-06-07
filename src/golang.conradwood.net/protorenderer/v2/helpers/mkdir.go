package helpers

import (
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
)

func Mkdir(dir string) error {
	err := linux.CreateIfNotExists(dir, 0777)
	if err != nil {
		fmt.Printf("failed to create dir \"%s\": %s\n", dir, err)
	} else {
		fmt.Printf("Created dir \"%s\"\n", dir)
	}
	return err

}
