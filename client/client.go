package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Moekr/sword/common"
	"github.com/Moekr/sword/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var args *util.Args

func Start(clientArgs *util.Args) error {
	args = clientArgs
	pingLoop()
	return nil
}

func pingLoop() {
	for {
		now := time.Now()
		next := time.Unix(0, (now.UnixNano()/int64(time.Minute)+1)*int64(time.Minute))
		time.Sleep(next.Sub(now))
		ping()
	}
}

func ping() (err error) {
	defer func() {
		if err != nil {
			log.Println("ping error: " + err.Error())
		}
	}()
	req, err := http.NewRequest(http.MethodGet, args.Server+"/api/conf", nil)
	if err != nil {
		return err
	}
	req.Header.Add(common.TokenHeaderName, args.Token)
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
			log.Println("ping target error: " + err.Error())
		}
	}()
	record := doPing(target.Address)
	bs, err := json.Marshal(record)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/data?t=%d&o=%d", args.Server, target.Id, args.ClientId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
	if err != nil {
		return err
	}
	req.Header.Add(common.TokenHeaderName, args.Token)
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
