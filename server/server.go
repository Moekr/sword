package server

import (
	"encoding/json"
	"github.com/Moekr/sword/util/args"
	"github.com/Moekr/sword/util/logs"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	_args    *args.Args
	conf     *Conf
	dataSets map[int64]map[int64]*DataSet
)

func Start(serverArgs *args.Args) error {
	_args = serverArgs
	if err := loadConf(); err != nil {
		return err
	}
	loadData()
	defer saveData(false)
	go refreshLoop()
	go deferKill()
	http.HandleFunc("/api/conf", httpConf)
	http.HandleFunc("/api/data", httpData)
	http.HandleFunc("/api/data/abbr", httpAbbrData)
	http.HandleFunc("/api/data/full", httpFullData)
	http.HandleFunc("/index.html", httpIndex)
	http.HandleFunc("/detail.html", httpDetail)
	http.HandleFunc("/static/index.css", httpCSS)
	http.HandleFunc("/static/index.js", httpJS)
	http.HandleFunc("/favicon.ico", httpFavicon)
	http.HandleFunc("/", httpIndex)
	return http.ListenAndServe(_args.Bind, nil)
}

func loadConf() error {
	if bs, err := ioutil.ReadFile(_args.ConfFile); err != nil {
		return err
	} else {
		return json.Unmarshal(bs, &conf)
	}
}

func refreshLoop() {
	for {
		now := time.Now()
		cur := time.Unix(0, now.UnixNano()-now.UnixNano()%int64(time.Minute))
		next := cur.Add(time.Minute)
		time.Sleep(next.Sub(now))
		for _, dataSets := range dataSets {
			for _, dataSet := range dataSets {
				dataSet.Refresh(cur)
			}
		}
		saveData(next.Minute() == 0)
	}
}

func deferKill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	sig := <-ch
	logs.Warn("receive signal %v", sig)
	saveData(false)
	if sig == syscall.SIGTERM {
		os.Exit(1)
	}
	os.Exit(0)
}
