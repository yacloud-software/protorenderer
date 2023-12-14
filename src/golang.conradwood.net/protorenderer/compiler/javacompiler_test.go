package compiler

import (
	"testing"
)

func TestPackagenameToJavaDir(t *testing.T) {
	res := make(map[string]string)
	res["golang.singingcat.net/apis"] = "net/singingcat/golang/apis"
	res["golang.singingcat.net"] = "net/singingcat/golang"
	res["golang.singingcat.net/apis/common"] = "net/singingcat/golang/apis/common"
	for k, v := range res {
		s := goPackagenameToJavaDir(k)
		if s != v {
			t.Errorf("%s: expected %s got %s", k, v, s)
		}
	}
}







































































































