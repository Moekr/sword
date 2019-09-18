package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Moekr/gopkg/logs"
	"github.com/Moekr/sword/server/conf"
	"github.com/Moekr/sword/server/cronjob"
	"github.com/Moekr/sword/server/dataset"
	"github.com/Moekr/sword/server/http"
	"github.com/Moekr/sword/server/persistence"
)

func Start() error {
	if err := conf.InitConf(); err != nil {
		return err
	}
	dataset.InitDataSets()
	if bs, err := persistence.LoadData(); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		dataset.Decode(bs)
	}
	http.StartHTTPService()
	cronjob.StartCronJob()
	return waitForSignal()
}

func waitForSignal() error {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	logs.Warn("[Server] received signal %v, save data and exit...", <-ch)
	return persistence.StoreData(false)
}
