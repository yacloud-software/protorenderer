package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/rpc"
)

func NeedVersion(ctx context.Context) error {
	if completeVersion != nil {
		return nil
	}
	cs := rpc.CallStateFromContext(ctx)
	s := fmt.Sprintf(fmt.Sprintf("%s.%s not available", cs.ServiceName, cs.MethodName))
	fmt.Println(s)
	return errors.Unavailable(ctx, s)
}
