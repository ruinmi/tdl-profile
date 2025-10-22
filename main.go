package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/go-faster/errors"
	"github.com/gotd/td/tg"
	"github.com/iyear/tdl/extension"
	"github.com/mattn/go-colorable"
	"github.com/neilotoole/jsoncolor"
	"github.com/spf13/pflag"
)

type profileUpdate struct {
	FirstName *string
	LastName  *string
	About     *string
}

func (u profileUpdate) Empty() bool {
	return u.FirstName == nil && u.LastName == nil && u.About == nil
}

func main() {
	update, showHelp, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if showHelp {
		return
	}

	extension.New(extension.Options{
		UpdateHandler: nil,
		Middlewares:   nil,
		Logger:        nil,
	})(func(ctx context.Context, e *extension.Extension) error {
		api := e.Client().API()
		if update.Empty() {
			full, err := api.UsersGetFullUser(ctx, &tg.InputUserSelf{})
			if err != nil {
				return errors.Wrap(err, "get current profile")
			}

			// 简化输出
			user := full.Users[0].(*tg.User)
			info := map[string]any{
				"id":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"username":  user.Username,
				"phone":     user.Phone,
				"about":     full.FullUser.About,
			}
		
			output("当前账号资料：", info, "")
			return nil
		}

		req := &tg.AccountUpdateProfileRequest{}
		if update.FirstName != nil {
			req.SetFirstName(*update.FirstName)
		}
		if update.LastName != nil {
			req.SetLastName(*update.LastName)
		}
		if update.About != nil {
			req.SetAbout(*update.About)
		}
		_, err := api.AccountUpdateProfile(ctx, req)
		if err != nil {
			return errors.Wrap(err, "update profile")
		}
		color.Green("资料更新成功！")
		return nil
	})
}

func parseFlags(args []string) (profileUpdate, bool, error) {
	var update profileUpdate

	var (
		firstName string
		lastName  string
		about     string
	)

	flagSet := pflag.NewFlagSet("tdl-profile", pflag.ContinueOnError)
	writer := colorable.NewColorable(os.Stdout)
	flagSet.SetOutput(writer)
	flagSet.StringVarP(&firstName, "first-name", "f", "", "新的名字")
	flagSet.StringVarP(&lastName, "last-name", "l", "", "新的姓氏")
	flagSet.StringVarP(&about, "about", "a", "", "新的自我介绍")
	flagSet.Usage = func() {
		fmt.Fprintln(writer, "Usage of tdl-profile:")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		if err == pflag.ErrHelp {
			return update, true, nil
		}
		return update, false, err
	}

	if flagSet.Lookup("first-name").Changed {
		update.FirstName = &firstName
	}
	if flagSet.Lookup("last-name").Changed {
		update.LastName = &lastName
	}
	if flagSet.Lookup("about").Changed {
		update.About = &about
	}

	return update, false, nil
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
