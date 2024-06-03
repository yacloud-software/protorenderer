package golang

import (
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/protorenderer/v2/config"
	"path/filepath"
	"strings"
)

func FindCompiler(cname string) string {
	check := []string{
		"dist/linux/amd64/",
		fmt.Sprintf("extra/compilers/%s/", config.GetCompilerVersion()),
		"linux/amd64/",
		cmdline.GetYACloudDir() + "/ctools/dev/go/current/go/bin/",
		"/opt/cnw/ctools/dev/go/current/go/bin/",
		"/home/cnw/go/bin/",
		"/home/cnw/devel/go/protorenderer/dist/linux/amd64/",
	}
	var err error
	for _, d := range check {
		c := d + cname
		if !strings.HasPrefix(d, "/") {
			c, err = utils.FindFile(d + cname)
			if err != nil {
				continue
			}
		}

		if !utils.FileExists(c) {
			continue
		}
		if c[0] == '/' {
			return c
		}
		cs, err := filepath.Abs(c)
		if err != nil {
			fmt.Printf("Unable to absolutise \"%s\": %s\n", cs, err)
			return c
		}
		return cs
	}
	e := fmt.Sprintf("%s not found\n", cname)
	fmt.Println(e)
	panic(e)
}
