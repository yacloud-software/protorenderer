package server

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	//plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	//	pr "golang.conradwood.net/apis/protorenderer"
	//	"path/filepath"
	"strings"
)

type CommentResult struct {
	Messages map[string]string
	Services map[string]string
	RPCs     map[string]string // key: "Servicename.Fieldname"
	Fields   map[string]string // key: "Messagename.Fieldname"
}

func (smc *ServerMetaCompiler) dealWithComments(pf *descriptor.FileDescriptorProto) (*CommentResult, error) {
	smc.debugf("----comments:\n")
	res := &CommentResult{
		Messages: make(map[string]string),
		Services: make(map[string]string),
		RPCs:     make(map[string]string),
		Fields:   make(map[string]string),
	}
	if pf.SourceCodeInfo == nil {
		// no source code info
		return res, nil
	}

	//	fqdn := filepath.Dir(*pf.Name)
	//	pkg := *pf.Package
	//	pkg := result.ObtainPackage(fqdn)

	for si, svc := range pf.Service {
		//		sv := pkg.ObtainService(*svc.Name)
		comment := FindCommentPath(pf, 6, si)
		comment_text := comment.CombinedText()
		smc.debugf(" name=%s, si=%d (%s)\n", *svc.Name, si, comment_text)
		res.Services[*svc.Name] = comment_text
	}

	// parse all comments for RPCs
	for si, svc := range pf.Service {
		//sv := pkg.ObtainService(*svc.Name)
		for mi, mp := range svc.Method {
			//	rpc := sv.ObtainRPC(*mp.Name)
			comment := FindCommentPath(pf, 6, si, 2, mi)
			comment_text := comment.CombinedText()
			smc.debugf("    %s: %s\n", *mp.Name, comment_text)
			key := fmt.Sprintf("%s.%s", *svc.Name, *mp.Name)
			res.RPCs[key] = comment_text
		}
	}

	// parse all comments for Messages & Fields
	for mi, m := range pf.MessageType {
		//		msg := pkg.ObtainMessage(*m.Name)
		comment := FindCommentPath(pf, 4, mi)
		comment_text := comment.CombinedText()
		smc.debugf("Comment for %s: %s\n", *m.Name, comment_text)
		res.Messages[*m.Name] = comment_text

		for _, f := range m.Field {
			// this really is very very confusing and I have no idea why i have to subtract 1 here
			comment := FindCommentPath(pf, 4, mi, 2, int(*f.Number)-1)
			//	fl := msg.ObtainField(*f.Name)
			comment_text := comment.CombinedText()
			smc.debugf("Comment for %s: %s\n", *f.Name, comment_text)
			field_name := fmt.Sprintf("%s.%s", *m.Name, *f.Name)
			res.Fields[field_name] = comment_text
		}
	}

	return res, nil
}

// comment by path...
func FindCommentPath(fdp *descriptor.FileDescriptorProto, p ...int) *Comment {
	for _, sciLoc := range fdp.SourceCodeInfo.Location {
		if isPath(sciLoc, p...) {
			return commentFromSciloc(sciLoc)
		}
	}
	return nil
}
func isPath(sciLoc *descriptor.SourceCodeInfo_Location, p ...int) bool {
	if len(sciLoc.Path) != len(p) {
		return false
	}
	for i, lp := range sciLoc.Path {
		if lp != int32(p[i]) {
			return false
		}
	}
	return true
}

type Comment struct {
	path     string
	leading  string
	trailing string
}

func (c *Comment) CombinedText() string {
	if c == nil {
		return ""
	}
	s := c.leading + c.trailing
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r")
	return s
}
func commentFromSciloc(sciLoc *descriptor.SourceCodeInfo_Location) *Comment {
	c := &Comment{}
	pname := ""
	for _, p := range sciLoc.Path {
		pname = pname + fmt.Sprintf("%d", p) + "_"
	}
	if len(sciLoc.Path) == 0 {
		pname = "empty"
	}
	c.path = pname
	if sciLoc.LeadingComments != nil {
		c.leading = *sciLoc.LeadingComments
	}
	if sciLoc.TrailingComments != nil {
		c.trailing = *sciLoc.TrailingComments
	}
	return c
}

func (c *CommentResult) GetMessageComment(name string) string {
	com, found := c.Messages[name]
	if !found {
		return ""
	}
	return com
}
func (c *CommentResult) GetFieldComment(messagename, fieldname string) string {
	key := fmt.Sprintf("%s.%s", messagename, fieldname)
	com, found := c.Fields[key]
	if !found {
		fmt.Printf("WARNING - no comment for \"%s\"\n", key)
		return ""
	}
	return com
}
