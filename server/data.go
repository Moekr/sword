package server

import (
	"github.com/Moekr/sword/common"
	"sync"
	"time"
)

const (
	timeIdx = iota
	avgIdx
	maxIdx
	minIdx
	lostIdx
	length
)

type DataSet struct {
	Target    *common.Target   `json:"target"`
	Observer  *common.Observer `json:"observer"`
	DayData   [][]int64        `json:"day_data"`
	WeekData  [][]int64        `json:"week_data"`
	MonthData [][]int64        `json:"month_data"`
	YearData  [][]int64        `json:"year_data"`
	Lock      *sync.RWMutex    `json:"-"`
}

type AbbrDataSet struct {
	Observer *common.Observer `json:"observer"`
	Data     []*AbbrData      `json:"data"`
}

type AbbrData struct {
	Time  int64 `json:"time"`
	Value int64 `json:"value"`
}

func NewEmptyDataSet(t *common.Target, o *common.Observer) *DataSet {
	dataSet := &DataSet{
		Target:    t,
		Observer:  o,
		DayData:   make([][]int64, 0),
		WeekData:  make([][]int64, 0),
		MonthData: make([][]int64, 0),
		YearData:  make([][]int64, 0),
	}
	dataSet.Init()
	return dataSet
}

func (d *DataSet) Init() {
	d.Lock = &sync.RWMutex{}
	d.initDayData()
}

func (d *DataSet) initDayData() {
	now := time.Now()
	cur := time.Unix(0, now.UnixNano()-now.UnixNano()%int64(time.Minute))
	fst := cur.Add(-1440 * time.Minute)
	oldData := make([][]int64, 0)
	for idx, data := range d.DayData {
		if data[timeIdx] == fst.Unix() {
			oldData = d.DayData[idx:]
		}
	}
	newData := make([][]int64, len(oldData), 1440)
	for idx, data := range oldData {
		newData[idx] = data
	}
	for i := len(newData); i < 1440; i++ {
		newData = append(newData, []int64{fst.Add(time.Duration(i) * time.Minute).Unix(), -1, -1, -1, -1})
	}
	d.DayData = newData
	go checkDayDataLoop(d)
}

func checkDayDataLoop(dataSet *DataSet) {
	for {
		checkDayData(dataSet)
	}
}

func checkDayData(dataSet *DataSet) {
	now := time.Now()
	cur := time.Unix(0, now.UnixNano()-now.UnixNano()%int64(time.Minute))
	next := time.Unix(0, (now.UnixNano()/int64(time.Minute)+1)*int64(time.Minute))
	time.Sleep(next.Sub(now))
	dataSet.Lock.Lock()
	defer dataSet.Lock.Unlock()
	offset := (cur.Unix() - dataSet.DayData[1439][timeIdx]) / 60
	for offset > 0 {
		dataSet.put(&common.Record{
			Time: next.Unix() - offset*60,
			Avg:  -1,
			Max:  -1,
			Min:  -1,
			Lost: -1,
		}, true)
		offset--
	}
}

func (d *DataSet) Put(record *common.Record) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.put(record, false)
}

func (d *DataSet) put(record *common.Record, noCheck bool) {
	now := time.Now()
	cur := time.Unix(0, now.UnixNano()-now.UnixNano()%int64(time.Minute))
	if noCheck || cur.Unix() > d.DayData[1439][timeIdx] {
		for i := 0; i < 1439; i++ {
			d.DayData[i] = d.DayData[i+1]
		}
		d.DayData[1439] = []int64{cur.Unix(), record.Avg, record.Max, record.Min, record.Lost}
	}
}

func (d *DataSet) GetAbbrData() *AbbrDataSet {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	data := make([]*AbbrData, len(d.DayData))
	for idx, record := range d.DayData {
		data[idx] = &AbbrData{
			Time:  record[timeIdx],
			Value: record[avgIdx],
		}
	}
	return &AbbrDataSet{
		Observer: d.Observer,
		Data:     data,
	}
}
