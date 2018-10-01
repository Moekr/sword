package server

import (
	"encoding/json"
	"fmt"
	"github.com/Moekr/sword/util"
	"io/ioutil"
	"path"
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

func saveData() {
	d := args.DataDir
	util.Infof("save data to %s\n", d)
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
