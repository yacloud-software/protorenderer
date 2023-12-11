package meta

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pr "golang.conradwood.net/apis/protorenderer"
)

type Result struct {
	VerifyToken string
	Packages    []*Package
}

/************************************************************************
* Package
************************************************************************/

type Package struct {
	Name       string
	FQDN       string
	Proto      *pr.Package
	Filename   string // which file created this
	Comment    string
	Services   []*Service
	Messages   []*Message
	Protofiles []*pr.ProtoFile
	CNWOptions map[string]string // options in comments, such as //CNW_OPTION: foo=bar
}

// give, e.g. "golang.conradwood.net/apis/common"
func (r *Result) ObtainPackage(fqdn string) *Package {
	for _, p := range r.Packages {
		if p.FQDN == fqdn {
			return p
		}
	}
	p := &Package{
		FQDN: fqdn,
	}
	r.Packages = append(r.Packages, p)
	return p
}

/************************************************************************
* Service
************************************************************************/
type Service struct {
	ID      string
	Name    string
	Comment string
	RPCs    []*RPC
}

func (p *Package) ObtainService(name string) *Service {
	for _, s := range p.Services {
		if s.Name == name {
			return s
		}
	}
	s := &Service{Name: name}
	p.Services = append(p.Services, s)
	return s
}

/************************************************************************
* RPC
************************************************************************/
type RPC struct {
	ID      string
	Name    string
	Comment string
	Service *Service
	Input   *Message
	Output  *Message
}

func (s *Service) ObtainRPC(name string) *RPC {
	for _, r := range s.RPCs {
		if r.Name == name {
			return r
		}
	}
	r := &RPC{Name: name, Service: s}
	s.RPCs = append(s.RPCs, r)
	return r
}

/************************************************************************
* Message
************************************************************************/
type Message struct {
	ID      string
	Name    string
	Comment string
	Package *Package
	Fields  []*Field
}

func (p *Package) ObtainMessage(name string) *Message {
	for _, m := range p.Messages {
		if m.Name == name {
			return m
		}
	}
	m := &Message{Name: name, Package: p}
	p.Messages = append(p.Messages, m)
	return m
}

func (m *Message) ToProto() *pr.Message {
	return nil
}

/************************************************************************
* Field
************************************************************************/
type Field struct {
	ID       string
	Name     string
	Comment  string
	Repeated bool
	Optional bool
	Required bool
	Type     descriptor.FieldDescriptorProto_Type
	Message  *Message
}

func (m *Message) ObtainField(name string) *Field {
	for _, f := range m.Fields {
		if f.Name == name {
			return f
		}
	}
	f := &Field{Name: name}
	m.Fields = append(m.Fields, f)
	return f
}

func (f *Field) TypeName() string {
	return fmt.Sprintf("%v", f.Type)
}








