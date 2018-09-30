package util

import "flag"

type Args struct {
	// for server
	IsServer bool   // -s
	Bind     string // -b
	ConfFile string // -c
	DataDir  string // -d
	// for client
	Server   string // -u
	ClientId int64  // -i
	// for both
	Token string // -t
}

func ParseArgs() *Args {
	args := &Args{}
	flag.BoolVar(&args.IsServer, "s", false, "identify server role or not")
	flag.StringVar(&args.Bind, "b", "0.0.0.0:7901", "address and port to bind")
	flag.StringVar(&args.ConfFile, "c", "./conf.json", "config file containing targets info")
	flag.StringVar(&args.DataDir, "d", "./data/", "data directory")
	flag.StringVar(&args.Server, "u", "http://localhost:7901", "server address")
	flag.Int64Var(&args.ClientId, "i", -1, "client id")
	flag.StringVar(&args.Token, "t", "", "token used in communication between server and client")
	flag.Parse()
	return args
}