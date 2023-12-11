package meta

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
	//	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	pr "golang.conradwood.net/apis/protorenderer"
	"path/filepath"
)

// this is called by protoc plugin (via SubmitSource)
func (m *MetaCompiler) generate(req *pr.ProtocRequest) error {
	pp := m.processingProtofile.Protofile()
	if pp == nil {
		return fmt.Errorf("missing protofile")
	}
	result := m.result
	if result == nil {
		result = &Result{}
	}
	var msgs []*Message
	// run through all messages & base-type fields
	for _, pf := range req.ProtoFiles {
		/*
			if pf.Syntax == nil {
				return fmt.Errorf("proto version not specified (in %s)", *pf.Name)
			}
			if *pf.Syntax != "proto3" {
				return fmt.Errorf("Only proto3 is supported (in %s: \"%s\")", *pf.Name, *pf.Syntax)
			}
		*/
		opts := parseCNWOptionsFromFile(pf)
		fqdn := filepath.Dir(*pf.Name)
		pkg := result.ObtainPackage(fqdn)
		pkg.CNWOptions = opts
		pkg.Protofiles = append(pkg.Protofiles, pp)
		pkg.Filename = *pf.Name
		pkg.Name = *pf.Package
		// TODO: check for discrepance in proto.Name
		if pkg.Proto == nil {
			pkg.Proto = &pr.Package{Name: *pf.Name, Prefix: fqdn}
		}
		debugf("   pkg=%s, name=%s, fqdn=%s\n", *pf.Name, *pf.Package, fqdn)

		// get all Messages
		for _, ms := range pf.MessageType {
			debugf("   Message: %s\n", *ms.Name)
			msg := pkg.ObtainMessage(*ms.Name)
			msgs = append(msgs, msg)
			for _, fd := range ms.Field {
				if fd.Type != nil && *fd.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					// in "Pass 1", we only deal with simple types not sub messages
					continue
				}
				fi := msg.ObtainField(*fd.Name)
				debugf("       Field: %s\n", *fd.Name)
				debugf("            Label   : %s\n", *fd.Label) // LABEL_OPTIONAL//LABEL_REPEATED
				if fd.TypeName != nil {
					debugf("            Typename: %s\n", *fd.TypeName)
				}
				if fd.Type != nil {
					debugf("            Type    : %v\n", *fd.Type)
					fi.Type = *fd.Type
				}
			}
		}
	}

	// run through all copmlex fields
	debugf("Complex types..\n")
	for _, pf := range req.ProtoFiles {
		fqdn := filepath.Dir(*pf.Name)
		pkg := result.ObtainPackage(fqdn)
		// get all Messages
		for _, ms := range pf.MessageType {
			debugf("   Message: %s\n", *ms.Name)
			msg := pkg.ObtainMessage(*ms.Name)
			for _, fd := range ms.Field {
				if fd.Type == nil || *fd.Type != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					// in "Pass 2", we only deal with complex types
					continue
				}
				fi := msg.ObtainField(*fd.Name)
				if fd.Label != nil {
					lf := *fd.Label
					if lf == descriptor.FieldDescriptorProto_LABEL_REPEATED {
						fi.Repeated = true
					} else if lf == descriptor.FieldDescriptorProto_LABEL_REQUIRED {
						fi.Required = true
					} else if lf == descriptor.FieldDescriptorProto_LABEL_OPTIONAL {
						fi.Optional = true
					}
				}
				fi.Type = *fd.Type
				debugf("       Field: %s\n", *fd.Name)
				debugf("            Label   : %s\n", *fd.Label) // LABEL_OPTIONAL//LABEL_REPEATED
				if fd.TypeName != nil {
					debugf("            Typename: %s\n", *fd.TypeName)
				}
				if fd.Type != nil {
					debugf("            Type    : %v\n", *fd.Type)
				}
				fi.Message = ResolveMessage(msgs, *fd.TypeName)
			}
		}
	}

	// create all RPCs
	for _, pf := range req.ProtoFiles {
		fqdn := filepath.Dir(*pf.Name)
		pkg := result.ObtainPackage(fqdn)
		for _, sv := range pf.Service {
			debugf("   Service: %s\n", *sv.Name)
			svc := pkg.ObtainService(*sv.Name)
			for _, mp := range sv.Method {
				rpc := svc.ObtainRPC(*mp.Name)
				rpc.Input = ResolveMessage(msgs, *mp.InputType)
				rpc.Output = ResolveMessage(msgs, *mp.OutputType)
				//				fmt.Printf("       RPC: %s(%s) %s\n", *mp.Name, *mp.InputType, *mp.OutputType)
			}
		}
	}

	err := dealWithComments(req, result)
	if err != nil {
		return err
	}
	debugf("\n\n--------------- result -------------\n\n")

	err = m.submitResult(result) // assign persistent IDs
	if err != nil {
		fmt.Printf("Failed to meta compile: %s\n", err)
		return err
	}
	m.result = result
	return nil
}

// given a list of messages and a fqdn, like '.common.void', it will return the message
// since the exact rules of fqdns aren't clear, it'll ONLY resolve .[package].[name],
// if it is something like .[something].[package].[name] - it will return nil
func ResolveMessage(msgs []*Message, fqdn string) *Message {
	if len(fqdn) == 0 {
		s := fmt.Sprintf("Attempt to resolve message with empty fqdn")
		fmt.Println(s)
		return nil
	}
	if fqdn[0] != '.' {
		s := fmt.Sprintf("Attempt to resolve message without fqdn (does not start with dot): \"%s\"", fqdn)
		fmt.Println(s)
		return nil
	}
	pkgname := ".google.protobuf."
	msgname := ""
	if strings.HasPrefix(fqdn, pkgname) {
		i := strings.LastIndex(fqdn, ".")
		pkgname = strings.Trim(fqdn[1:i], ".")
		msgname = strings.Trim(fqdn[i:], ".")

	} else {
		fqdn = fqdn[1:]
		sx := strings.Split(fqdn, ".")
		if len(sx) != 2 {
			s := fmt.Sprintf("Unhandled message fqdn (expected 2 parts, not %d): %s", len(sx), fqdn)
			fmt.Println(s)
			return nil
		}
		pkgname = sx[0]
		msgname = sx[1]
	}
	for _, m := range msgs {
		p := m.Package
		if p.Name != pkgname || m.Name != msgname {
			continue
		}
		return m
	}
	if !strings.HasPrefix(fqdn, ".google.protobuf") {
		fmt.Printf("WARNING - did not find message for \"%s\" (pkg=%s,msg=%s)\n", fqdn, pkgname, msgname)
	}
	return nil
}





































































