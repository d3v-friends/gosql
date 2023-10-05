package main

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnLogger"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/d3v-friends/gosql"
	"github.com/spf13/cobra"
)

func CmdReader() (res *cobra.Command) {
	res = &cobra.Command{}
	res.Use = "reader"

	var fPath = "path"
	res.Flags().String(fPath, "./gosql.yaml", "--path ./config.yaml")

	res.Run = func(cmd *cobra.Command, args []string) {
		var logger = fnLogger.NewDefaultLogger()
		logger.SetLevel(fnLogger.Trace)

		var path = fnPanic.OnValue(cmd.Flags().GetString(fPath))

		logger.Trace("path: %s", path)
		fnPanic.IsTrue(path != "", fmt.Errorf("invalid path: path=%s", path))

		var config = fnPanic.OnPointer(gosql.Read(gosql.NewPath(path)))
		logger.Trace("config loaded")
		fnPanic.On(config.Validate())

		logger.Trace("config validate success")

	}

	return
}
