package main

import (
	"context"
	"fmt"
	"github.com/d3v-friends/go-pure/fnLogger"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/gosql"
	"github.com/d3v-friends/gosql/gctl"
	"github.com/d3v-friends/gosql/gdriver"
	"github.com/spf13/cobra"
)

func CmdMigrate() (res *cobra.Command) {

	res = &cobra.Command{}
	res.Use = "migrate"

	var fPath = "path"
	res.Flags().String(fPath, "./gosql.yaml", "--path ./config.yaml")

	var fUsername = "username"
	res.Flags().String(fUsername, "", "--username root")

	var fPassword = "password"
	res.Flags().String(fPassword, "", "--password 1234")

	var fHost = "host"
	res.Flags().String(fHost, "", "--host 123.123.123:3306")

	res.Run = func(cmd *cobra.Command, args []string) {
		var logger = fnLogger.NewDefaultLogger()
		logger.SetLevel(fnLogger.Trace)

		var path = fnPanic.OnValue(cmd.Flags().GetString(fPath))
		var username = fnPanic.OnValue(cmd.Flags().GetString(fUsername))
		var password = fnPanic.OnValue(cmd.Flags().GetString(fPassword))
		var host = fnPanic.OnValue(cmd.Flags().GetString(fHost))

		logger.Trace("path: %s", path)
		fnPanic.IsTrue(path != "", fmt.Errorf("invalid path: path=%s", path))

		var cfg = fnPanic.OnPointer(gosql.Read(gosql.NewPath(path)))
		logger.Trace("config loaded")
		fnPanic.On(cfg.Validate())

		logger.Trace("config validate success")

		// migrate 실행
		var driver, db, err = gdriver.NewMySQL5(cfg.Type, &gdriver.IConn{
			Username: username,
			Password: password,
			Host:     host,
		})

		if err != nil {
			panic(err)
		}

		var ctx = context.TODO()
		ctx = fnLogger.Set(ctx, logger)
		if err = gctl.Migrate(ctx, db, driver, cfg); err != nil {
			panic(err)
		}
	}

	return
}
