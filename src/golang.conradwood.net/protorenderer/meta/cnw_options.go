package meta

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
)

const (
	CNW_OPTION_KEYWORD = "CNW_OPTION: "
)

func parseCNWOptionsFromFile(fdp *descriptor.FileDescriptorProto) map[string]string {
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
			handle_cnw_option_string(res, lc)
		}
		if l.TrailingComments != nil {
			tc = *l.TrailingComments
			handle_cnw_option_string(res, tc)
		}
		//		debugf("Leading Comments: \"%s\"\n", lc)
		//		debugf("Trailing Comments: \"%s\"\n", tc)
		for _, ldc := range l.LeadingDetachedComments {
			handle_cnw_option_string(res, ldc)
		}
	}
	return res
}
func loc_to_string(l []int32) string {
	return fmt.Sprintf("%v", l)

}
func handle_cnw_option_string(res map[string]string, opt string) {
	if !strings.Contains(opt, CNW_OPTION_KEYWORD) {
		return
	}
	opt = strings.TrimPrefix(opt, " ")
	opt = strings.TrimSuffix(opt, "\n")
	if !strings.HasPrefix(opt, CNW_OPTION_KEYWORD) {
		return
	}
	opt = strings.TrimPrefix(opt, CNW_OPTION_KEYWORD)
	debugf("Cnw option: \"%s\"\n", opt)
}




























































































