package main

type StandardProtoFile struct {
	id       uint64
	filename string
	content  []byte
}

func (spf *StandardProtoFile) GetID() uint64 {
	return spf.id
}
func (spf *StandardProtoFile) GetFilename() string {
	return spf.filename
}
func (spf *StandardProtoFile) Content() []byte {
	return spf.content
}
