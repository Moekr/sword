package main

import (
	"github.com/Moekr/sword/client"
	"github.com/Moekr/sword/server"
	"github.com/Moekr/sword/util/args"
	"github.com/Moekr/sword/util/logs"
)

func main() {
	_args := args.Parse()
	logs.SetDebug(_args.IsDebug)
	var err error
	if _args.IsServer {
		err = server.Start(_args)
	} else {
		err = client.Start(_args)
	}
	logs.Fatal("exit with error: %s", err.Error())
}
