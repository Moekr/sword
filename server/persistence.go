package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

func loadData() {
	d := args.DataDir
	log.Println("load data from " + d)
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
				}
			}
			dataSets[target.Id][observer.Id] = NewEmptyDataSet(target, observer)
		}
	}
}

func saveData() {
	d := args.DataDir
	log.Println("save data to " + d)
	for _, target := range conf.Targets {
		for _, observer := range conf.Observers {
			f := fmt.Sprintf("sword-%d-%d.json", target.Id, observer.Id)
			if bs, err := json.Marshal(dataSets[target.Id][observer.Id]); err == nil {
				ioutil.WriteFile(path.Join(d, f), bs, 0755)
			}
		}
	}
}
