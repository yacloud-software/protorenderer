package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
)

func NeedVersion(ctx context.Context) error {
	if completeVersion != nil {
		return nil
	}
	svc := auth.UserIDString(auth.GetService(ctx))
	uvc := auth.UserIDString(auth.GetUser(ctx))
	s := fmt.Sprintf(fmt.Sprintf("%s.%s not available", uvc, svc))
	fmt.Println(s)
	return errors.Unavailable(ctx, s)
}








































