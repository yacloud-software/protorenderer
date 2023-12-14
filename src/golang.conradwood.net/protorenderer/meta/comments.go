package meta

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	//plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	pr "golang.conradwood.net/apis/protorenderer"
	"path/filepath"
	"strings"
)

func dealWithComments(req *pr.ProtocRequest, result *Result) error {
	debugf("----comments:\n")

	// parse all the comments for services
	for _, pf := range req.ProtoFiles {
		if pf.SourceCodeInfo == nil {
			// no source code info
			continue
		}
		fqdn := filepath.Dir(*pf.Name)
		pkg := result.ObtainPackage(fqdn)
		for si, svc := range pf.Service {
			sv := pkg.ObtainService(*svc.Name)
			comment := FindCommentPath(pf, 6, si)
			sv.Comment = comment.CombinedText()
			debugf(" name=%s, si=%d (%s)\n", *svc.Name, si, sv.Comment)
		}
	}

	// parse all comments for RPCs
	for _, pf := range req.ProtoFiles {
		if pf.SourceCodeInfo == nil {
			// no source code info
			continue
		}
		fqdn := filepath.Dir(*pf.Name)
		pkg := result.ObtainPackage(fqdn)
		for si, svc := range pf.Service {
			sv := pkg.ObtainService(*svc.Name)
			for mi, mp := range svc.Method {
				rpc := sv.ObtainRPC(*mp.Name)
				comment := FindCommentPath(pf, 6, si, 2, mi)
				rpc.Comment = comment.CombinedText()
				debugf("    %s: %s\n", rpc.Name, rpc.Comment)

			}
		}
	}

	// parse all comments for Messages & Fields
	for _, pf := range req.ProtoFiles {
		if pf.SourceCodeInfo == nil { // no source code info
			continue
		}
		fqdn := filepath.Dir(*pf.Name)
		pkg := result.ObtainPackage(fqdn)
		for mi, m := range pf.MessageType {
			msg := pkg.ObtainMessage(*m.Name)
			comment := FindCommentPath(pf, 4, mi)
			msg.Comment = comment.CombinedText()

			for _, f := range m.Field {
				// this really is very very confusing and I have no idea why i have to subtract 1 here
				comment := FindCommentPath(pf, 4, mi, 2, int(*f.Number)-1)
				fl := msg.ObtainField(*f.Name)
				fl.Comment = comment.CombinedText()
				//	debugf("Comment for %s: %s\n", *f.Name, fl.Comment)
			}
		}
	}

	return nil
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









































































































