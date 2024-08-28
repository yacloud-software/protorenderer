package srv

import (
	"context"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/protorenderer/v2/store"
)

func (pr *protoRenderer) SubmitStore(ctx context.Context, req *common.Void) (*common.Void, error) {
	store.TriggerUpload(CompileEnv.StoreDir())

	return req, nil
}
