package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/iyear/tdl/extension"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	conf := zap.NewDevelopmentConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	extension.New(extension.Options{
		UpdateHandler: nil,
		Middlewares:   nil,
		Logger: nil,
	})(func(ctx context.Context, e *extension.Extension) error {
		self, err := e.Client().Self(ctx)
		if err != nil {
			return errors.Wrap(err, "get self")
		}

		api := client.API()
		req := &tg.AccountUpdateProfileRequest{}
		req.SetAbout("我的新自我介绍 ✨")
		_, err := api.AccountUpdateProfile(ctx, req)
		if err != nil {
			return errors.Wrap(err, "update profile")
		}
		return nil
	})
}
