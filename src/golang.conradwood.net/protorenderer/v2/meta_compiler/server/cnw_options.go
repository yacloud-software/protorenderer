package server

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
)

const (
	CNW_OPTION_KEYWORD = "CNW_OPTION: "
)

func (smc *ServerMetaCompiler) parseCNWOptionsFromFile(fdp *descriptor.FileDescriptorProto) map[string]string {
	res := make(map[string]string)
	sci := fdp.SourceCodeInfo
	if sci == nil {
		return res
	}
	for _, l := range sci.Location {
		//		debugf("Location: %s\n", loc_to_string(l.Path))
		lc := ""
		tc := ""
		if l.LeadingComments != nil {
			lc = *l.LeadingComments
			smc.handle_cnw_option_string(res, lc)
		}
		if l.TrailingComments != nil {
			tc = *l.TrailingComments
			smc.handle_cnw_option_string(res, tc)
		}
		//		debugf("Leading Comments: \"%s\"\n", lc)
		//		debugf("Trailing Comments: \"%s\"\n", tc)
		for _, ldc := range l.LeadingDetachedComments {
			smc.handle_cnw_option_string(res, ldc)
		}
	}
	return res
}
func loc_to_string(l []int32) string {
	return fmt.Sprintf("%v", l)

}
func (smc *ServerMetaCompiler) handle_cnw_option_string(res map[string]string, opt string) {
	if !strings.Contains(opt, CNW_OPTION_KEYWORD) {
		return
	}
	opt = strings.TrimPrefix(opt, " ")
	opt = strings.TrimSuffix(opt, "\n")
	if !strings.HasPrefix(opt, CNW_OPTION_KEYWORD) {
		return
	}
	opt = strings.TrimPrefix(opt, CNW_OPTION_KEYWORD)
	smc.debugf("Cnw option: \"%s\"\n", opt)
	kv := strings.SplitN(opt, "=", 2)
	if len(kv) == 2 {
		res[kv[0]] = kv[1]

	} else {
		res[opt] = ""
	}

}
