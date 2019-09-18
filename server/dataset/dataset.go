package dataset

import (
	"sync"
	"time"

	"github.com/Moekr/sword/server/conf"
	"github.com/Moekr/sword/types"
)

const (
	weekSR  = 7 * time.Minute
	monthSR = 30 * time.Minute
	yearSR  = 6 * time.Hour
)

var (
	DataSets types.TDataSets
)

func InitDataSets() {
	targets, clients := conf.Conf.Targets, conf.Conf.Clients
	DataSets = make(map[int16]*types.TTDataSet, len(targets))
	for _, target := range targets {
		m := make(map[int16]*types.TCDataSet, len(clients))
		for _, client := range clients {
			m[client.ID] = &types.TCDataSet{
				Client: client,
				Data:   newDataSet(),
			}
		}
		DataSets[target.ID] = &types.TTDataSet{
			Target: target,
			Data:   m,
		}
	}
}

func newDataSet() *types.TDataSet {
	makeSlice := func() []*types.TRecord {
		slice := make([]*types.TRecord, 1440)
		for idx := range slice {
			slice[idx] = types.NewRecord()
		}
		return slice
	}
	return &types.TDataSet{
		Lock:  &sync.RWMutex{},
		Day:   makeSlice(),
		Week:  makeSlice(),
		Month: makeSlice(),
		Year:  makeSlice(),
	}
}

func UpdateDataSets() {
	for _, tds := range DataSets {
		for _, cds := range tds.Data {
			updateDataSet(cds.Data)
		}
	}
}

func updateDataSet(ds *types.TDataSet) {
	ds.Lock.Lock()
	defer ds.Lock.Unlock()
	ts := time.Now().UnixNano() / int64(time.Minute) * int64(time.Minute)
	for i := 0; i < 1439; i++ {
		ds.Day[i] = ds.Day[i+1]
	}
	if ds.Buf != nil {
		ds.Day[1439], ds.Buf = ds.Buf, nil
	} else {
		ds.Day[1439] = types.NewRecord()
	}
	if ts%int64(weekSR) == 0 {
		for i := 0; i < 1439; i++ {
			ds.Week[i] = ds.Week[i+1]
		}
		ds.Week[1439] = Union(ds.Day[1433:])
	}
	if ts%int64(monthSR) == 0 {
		for i := 0; i < 1439; i++ {
			ds.Month[i] = ds.Month[i+1]
		}
		ds.Month[1439] = Union(ds.Day[1410:])
	}
	if ts%int64(yearSR) == 0 {
		for i := 0; i < 1439; i++ {
			ds.Year[i] = ds.Year[i+1]
		}
		ds.Year[1439] = Union(ds.Day[1080:])
	}
}
