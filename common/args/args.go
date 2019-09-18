package args

import "flag"

type TArgs struct {
	// for server
	RunAsServer bool   // -s
	BindAddress string // -b
	ConfPath    string // -c
	DataPath    string // -d
	// for client
	ServerAddress string // -u
	ClientId      int64  // -i
	// for both
	LogsPath  string // -l
	Token     string // -t
	DebugMode bool   // -v
}

var (
	Args TArgs
)

func InitArgs() {
	flag.BoolVar(&Args.RunAsServer, "s", false, "run as server or not")
	flag.StringVar(&Args.BindAddress, "b", "0.0.0.0:7901", "address to bind")
	flag.StringVar(&Args.ConfPath, "c", "/etc/sword/conf.json", "path of config file")
	flag.StringVar(&Args.DataPath, "d", "/opt/sword/data/", "path of data directory")
	flag.StringVar(&Args.ServerAddress, "u", "http://localhost:7901", "server address")
	flag.Int64Var(&Args.ClientId, "i", -1, "client id")
	flag.StringVar(&Args.LogsPath, "l", "/opt/sword/logs", "path of logs directory")
	flag.StringVar(&Args.Token, "t", "", "token used in communication between server and client")
	flag.BoolVar(&Args.DebugMode, "v", false, "enable debug mode or not")
	flag.Parse()
}
