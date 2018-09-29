package main

import (
	"github.com/Moekr/sword/client"
	"github.com/Moekr/sword/server"
	"github.com/Moekr/sword/util"
	"log"
)

func main() {
	args := util.ParseArgs()
	var err error
	if args.IsServer {
		err = server.Start(args)
	} else {
		err = client.Start(args)
	}
	log.Fatalf("exit with error: %s\n", err.Error())
}
