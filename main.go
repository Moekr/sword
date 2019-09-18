package main

import (
	"github.com/Moekr/gopkg/logs"
	"github.com/Moekr/sword/client"
	"github.com/Moekr/sword/common/args"
	"github.com/Moekr/sword/common/version"
	"github.com/Moekr/sword/server"
)

func main() {
	args.InitArgs()
	logs.InitLogs(args.Args.LogsPath)

	logs.Info("[Main] sword version: %s", version.Version)
	var err error
	if args.Args.RunAsServer {
		logs.Info("[Main] run as server role")
		err = server.Start()
	} else {
		logs.Info("[Main] run as client role")
		err = client.Start()
	}
	if err != nil {
		logs.Fatal("[Main] exit with error: %s", err.Error())
	}
}
