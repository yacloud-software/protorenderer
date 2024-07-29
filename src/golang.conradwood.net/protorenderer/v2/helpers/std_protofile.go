package helpers

type StandardProtoFile struct {
	id       uint64
	Filename string
	content  []byte
}

func (spf *StandardProtoFile) GetID() uint64 {
	return spf.id
}
func (spf *StandardProtoFile) GetFilename() string {
	return spf.Filename
}
func (spf *StandardProtoFile) Content() []byte {
	return spf.content
}
