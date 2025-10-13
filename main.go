package main

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/iyear/tdl/extension"
)

func main() {
	extension.New(extension.Options{
		UpdateHandler: nil,
		Middlewares:   nil,
		Logger: nil,
	})(func(ctx context.Context, e *extension.Extension) error {

		api := e.Client().API()
		req := &tg.AccountUpdateProfileRequest{}
		req.SetAbout("我的新自我介绍 ✨")
		_, err := api.AccountUpdateProfile(ctx, req)
		if err != nil {
			return errors.Wrap(err, "update profile")
		}
		return nil
	})
}
