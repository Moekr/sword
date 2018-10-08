package server

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Moekr/sword/util"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

func loadData() {
	d := args.DataDir
	util.Infof("load data from %s\n", d)
	for _, target := range conf.Targets {
		for _, observer := range conf.Observers {
			f := fmt.Sprintf("sword-%d-%d.json", target.Id, observer.Id)
			if bs, err := ioutil.ReadFile(path.Join(d, f)); err == nil {
				dataSet := &DataSet{}
				if err := json.Unmarshal(bs, dataSet); err == nil {
					dataSet.Target = target
					dataSet.Observer = observer
					dataSet.Init()
					dataSets[target.Id][observer.Id] = dataSet
					continue
				} else {
					util.Infof("unmarshal data file %s error: %s\n", f, err.Error())
				}
			} else {
				util.Infof("read data file %s error: %s\n", f, err.Error())
			}
			dataSets[target.Id][observer.Id] = NewEmptyDataSet(target, observer)
		}
	}
}

var lock = &sync.Mutex{}

func saveData() {
	lock.Lock()
	defer lock.Unlock()
	d := args.DataDir
	util.Debugf("save data to %s\n", d)
	for _, target := range conf.Targets {
		for _, observer := range conf.Observers {
			f := fmt.Sprintf("sword-%d-%d.json", target.Id, observer.Id)
			if bs, err := json.Marshal(dataSets[target.Id][observer.Id]); err == nil {
				if err = ioutil.WriteFile(path.Join(d, f), bs, 0755); err != nil {
					util.Infof("write data file %s error: %s\n", f, err.Error())
				}
			} else {
				util.Infof("marshal data file %s error: %s\n", f, err.Error())
			}
		}
	}
}

func backupData(cur time.Time) {
	d := args.DataDir
	f := fmt.Sprintf("backup-%d.tar.gz", cur.Unix())
	zipFile, err := os.Create(path.Join(d, f))
	if err != nil {
		util.Infof("create backup file %s error: %s\n", f, err.Error())
		return
	}
	defer zipFile.Close()
	gw := gzip.NewWriter(zipFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != d {
			return filepath.SkipDir
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			util.Debugf("open file %s error: %s\n", path, err.Error())
			return nil
		}
		if header, err := tar.FileInfoHeader(info, ""); err != nil {
			util.Debugf("create tar header %s error%s\n", path, err.Error())
		} else {
			header.Name = filepath.Base(header.Name)
			if err := tw.WriteHeader(header); err != nil {
				util.Debugf("add tar header %s error%s\n", path, err.Error())
			} else if _, err := io.Copy(tw, file); err != nil {
				util.Debugf("write tar file %s error\n", path, err.Error())
			}
		}
		return nil
	})
}
