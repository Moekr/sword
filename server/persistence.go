package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Moekr/sword/util"
	"github.com/Moekr/sword/util/logs"
	"os"
	"path"
	"sync"
	"time"
)

const (
	tmpFileName  = "sword-%d.json.gz"
	dataFileName = "sword.json.gz"
)

var lock = &sync.Mutex{}

func loadData() {
	p := path.Join(_args.DataDir, dataFileName)
	if err := loadDataImpl(p); err != nil {
		logs.Error("load data error: %s", err.Error())
		dataSets = nil
	}
	dss := make(map[int64]map[int64]*DataSet, len(conf.Targets))
	for _, target := range conf.Targets {
		dsm := make(map[int64]*DataSet, len(conf.Observers))
		for _, observer := range conf.Observers {
			if ds, ok := dataSets[target.Id][observer.Id]; ok {
				ds.Target = target
				ds.Observer = observer
				ds.Init()
				dsm[observer.Id] = ds
			} else {
				dsm[observer.Id] = NewEmptyDataSet(target, observer)
			}
		}
		dss[target.Id] = dsm
	}
	dataSets = dss
}

func loadDataImpl(p string) error {
	logs.Info("load data from %s begin", p)
	file, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("open file %s error: %s", p, err.Error())
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("create gzip reader error: %s", err.Error())
	}
	defer reader.Close()
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&dataSets); err != nil {
		return fmt.Errorf("decode from file %s error: %s", p, err.Error())
	}
	logs.Info("load data from %s success", p)
	return nil
}

func saveData(doBackup bool) {
	lock.Lock()
	defer lock.Unlock()
	d := _args.DataDir
	f := fmt.Sprintf(tmpFileName, time.Now().Unix())
	p := path.Join(d, f)
	if err := saveDataImpl(p); err != nil {
		logs.Error("save data error: %s", err.Error())
	} else {
		var fn func(src, dst string) error
		if doBackup {
			fn = util.Copy
		} else {
			fn = os.Rename
		}
		if err = fn(p, path.Join(d, dataFileName)); err != nil {
			logs.Error("save data error: %s", err.Error())
		}
	}
}

func saveDataImpl(p string) error {
	logs.Debug("save data to %s begin", p)
	file, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("create file %s error: %s", p, err.Error())
	}
	defer file.Close()
	writer := gzip.NewWriter(file)
	defer writer.Close()
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(dataSets); err != nil {
		return fmt.Errorf("encode to file %s error: %s", p, err.Error())
	}
	logs.Debug("save data to %s success", p)
	return nil
}
