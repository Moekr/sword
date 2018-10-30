package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Moekr/sword/common"
	"github.com/Moekr/sword/util/args"
	"github.com/Moekr/sword/util/logs"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var _args *args.Args

func Start(clientArgs *args.Args) error {
	_args = clientArgs
	go deferKill()
	pingLoop()
	return nil
}

func deferKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	logs.Warn("receive signal %v", <-ch)
	time.Sleep(time.Second)
	os.Exit(1)
}

func pingLoop() {
	for {
		now := time.Now()
		next := time.Unix(0, (now.UnixNano()/int64(time.Minute)+1)*int64(time.Minute))
		time.Sleep(next.Sub(now))
		if err := ping(); err != nil {
			logs.Debug("ping error: %s", err.Error())
		}
	}
}

func ping() (err error) {
	req, err := http.NewRequest(http.MethodGet, _args.Server+"/api/conf", nil)
	if err != nil {
		return err
	}
	req.Header.Add(common.TokenHeaderName, _args.Token)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response status %d when request conf", rsp.StatusCode)
	}
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	var targets []*common.Target
	if err := json.Unmarshal(bs, &targets); err != nil {
		return err
	}
	for _, target := range targets {
		go pingTarget(target)
	}
	return nil
}

func pingTarget(target *common.Target) (err error) {
	defer func() {
		if err != nil {
			logs.Debug("ping target error: %s", err.Error())
		}
	}()
	record := doPing(target.Address)
	bs, err := json.Marshal(record)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/data?t=%d&o=%d", _args.Server, target.Id, _args.ClientId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
	if err != nil {
		return err
	}
	req.Header.Add(common.TokenHeaderName, _args.Token)
	req.Header.Add("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response status %d when post data", rsp.StatusCode)
	}
	return nil
}
