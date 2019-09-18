package types

import "sync"

type TCategory struct {
	ID   int16        `json:"id",yaml:"id"`
	Name string       `json:"name",yaml:"name"`
	Sub  []*TCategory `json:"sub",yaml:"sub"`
}

type TTarget struct {
	ID   int16   `json:"id",yaml:"id"`
	Name string  `json:"name",yaml:"name"`
	Note string  `json:"note",yaml:"note"`
	Addr string  `json:"addr",yaml:"addr"`
	Cid  []int16 `json:"cid",yaml:"cid"`
}

type TClient struct {
	ID   int16  `json:"id",yaml:"id"`
	Name string `json:"name",yaml:"name"`
}

type TConf struct {
	Index      int16        `json:"index",yaml:"index"`
	Categories []*TCategory `json:"categories",yaml:"categories"`
	Targets    []*TTarget   `json:"targets",yaml:"targets"`
	Clients    []*TClient   `json:"clients",yaml:"clients"`
}

type TRecord struct {
	Avg int16 `json:"avg"`
	Max int16 `json:"max"`
	Min int16 `json:"min"`
	Std int16 `json:"std"`
	Los int8  `json:"los"`
}

type TDataSet struct {
	Lock  *sync.RWMutex
	Buf   *TRecord
	Day   []*TRecord
	Week  []*TRecord
	Month []*TRecord
	Year  []*TRecord
}

type TCDataSet struct {
	Client *TClient
	Data   *TDataSet
}

type TTDataSet struct {
	Target *TTarget
	Data   map[int16]*TCDataSet
}

type TDataSets map[int16]*TTDataSet
