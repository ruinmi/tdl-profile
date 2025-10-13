package main

import (
	"context"
	"os"
	"fmt"
	"github.com/fatih/color"

	"github.com/go-faster/errors"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/extension"
	"github.com/neilotoole/jsoncolor"
	"github.com/mattn/go-colorable"
	"github.com/spf13/pflag"
)

var firstName string
var lastName string
var about string

func init() {
	pflag.StringVarP(&firstName, "first-name", "f", "!@#^%$@#$&$&*&@)*(_).,;oyfgdasvbn", "新的名字")
	pflag.StringVarP(&lastName, "last-name", "l", "!@#^%$@#$&$&*&@)*(_).,;oyfgdasvbn", "新的姓氏")
	pflag.StringVarP(&about, "about", "a", "!@#^%$@#$&$&*&@)*(_).,;oyfgdasvbn", "新的自我介绍")
	pflag.Parse()
}

func main() {
	extension.New(extension.Options{
		UpdateHandler: nil,
		Middlewares:   nil,
		Logger: nil,
	})(func(ctx context.Context, e *extension.Extension) error {
		if firstName == "" && lastName == "" && about == "" {
			output("没有提供任何更新内容", nil, "")
			return nil
		}

		api := e.Client().API()
		req := &tg.AccountUpdateProfileRequest{}
		if firstName != "!@#^%$@#$&$&*&@)*(_).,;oyfgdasvbn" {
			req.SetFirstName(firstName)
		}
		if lastName != "!@#^%$@#$&$&*&@)*(_).,;oyfgdasvbn" {
			req.SetLastName(lastName)
		}
		if about != "!@#^%$@#$&$&*&@)*(_).,;oyfgdasvbn" {
			req.SetAbout(about)
		}
		_, err := api.AccountUpdateProfile(ctx, req)
		if err != nil {
			return errors.Wrap(err, "update profile")
		}
		return nil
	})
}


func output(header string, v any, footer string) {
	if header != "" {
		color.Blue(header)
	}

	enc := jsoncolor.NewEncoder(colorable.NewColorable(os.Stdout))

	clrs := jsoncolor.DefaultColors()
	clrs.Key = []byte("\x1b[35m")  // magenta
	clrs.Bool = []byte("\x1b[33m") // yellow
	enc.SetColors(clrs)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		fmt.Printf("%+v\n", v) // fallback
	}

	if footer != "" {
		color.Blue(footer)
	}

	fmt.Println()
}