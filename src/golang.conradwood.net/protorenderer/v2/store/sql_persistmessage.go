package store

import (
	"context"
	pb2 "golang.conradwood.net/apis/protorenderer2"
	//	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/protorenderer/db"
)

func GetOrCreateMessage(ctx context.Context, fileid uint64, name string) (*pb2.SQLMessage, error) {
	protofile, err := db.DefaultDBDBProtoFile().ByID(ctx, fileid)
	if err != nil {
		return nil, err
	}

	msgs, err := db.DefaultDBSQLMessage().ByName(ctx, name)
	if err != nil {
		return nil, err
	}
	var msgs2 []*pb2.SQLMessage
	for _, m := range msgs {
		if m.ProtoFile == nil || m.ProtoFile.ID != fileid {
			continue
		}
		msgs2 = append(msgs2, m)
	}
	if len(msgs2) != 0 {
		return msgs2[0], nil
	}
	dd := &pb2.SQLMessage{ProtoFile: protofile, Name: name}
	_, err = db.DefaultDBSQLMessage().Save(ctx, dd)
	if err != nil {
		return nil, err
	}
	return dd, nil
}
