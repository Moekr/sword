package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Moekr/gopkg/logs"
	"github.com/Moekr/gopkg/periodic"
	"github.com/Moekr/sword/client/ping"
	"github.com/Moekr/sword/common/args"
	"github.com/Moekr/sword/common/constant"
	"github.com/Moekr/sword/types"
)

func Start() error {
	periodic.NewStaticPeriodic(doJob, time.Minute, periodic.FixedRate).Start()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	logs.Warn("[Client] received signal %v, save data and exit...", <-ch)
	return nil
}

func doJob() {
	req, _ := http.NewRequest(http.MethodGet, args.Args.ServerAddress+"/api/conf", nil)
	req.Header.Add(constant.TokenHeader, args.Args.Token)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error("[Client] fetch conf request error: %s", err.Error())
		return
	} else if rsp.StatusCode != http.StatusOK {
		logs.Error("[Client] fetch conf response status: %d", rsp.StatusCode)
		return
	}
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logs.Error("[Client] fetch conf read error: %s", err.Error())
		return
	}
	var m struct {
		Code int
		Data []*types.TTarget
	}
	if err := json.Unmarshal(bs, &m); err != nil {
		logs.Error("[Client] fetch conf unmarshal error: %s", err.Error())
		return
	}
	for _, target := range m.Data {
		go asyncWorker(target)
	}
}

func asyncWorker(target *types.TTarget) {
	rec := ping.Ping(target.Addr)
	bs, _ := json.Marshal(rec)
	u := fmt.Sprintf("%s/api/push?t=%d&c=%d", args.Args.ServerAddress, target.ID, args.Args.ClientId)
	req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(bs))
	req.Header.Add(constant.TokenHeader, args.Args.Token)
	req.Header.Add("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error("[Client] upload data error: %s", err.Error())
	} else if rsp.StatusCode != http.StatusOK {
		logs.Error("[Client] upload data status: %d", rsp.StatusCode)
	}
}
