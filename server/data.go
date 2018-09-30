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
	Buffer    *common.Record   `json:"buffer"`
	Lock      *sync.RWMutex    `json:"-"`
}

type AbbrDataSet struct {
	Observer *common.Observer `json:"observer"`
	Data     [][]int64        `json:"data"`
}

type FullDataSet struct {
	Observer *common.Observer `json:"observer"`
	Data     [][]int64        `json:"data"`
}

var initNow = time.Now()

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
	cur := time.Unix(0, initNow.UnixNano()-initNow.UnixNano()%int64(time.Minute))
	fst := cur.Add(-1440 * time.Minute)
	oldData := make([][]int64, 0)
	for idx, data := range d.DayData {
		if data[timeIdx] == fst.Unix() {
			oldData = d.DayData[idx:]
			break
		}
	}
	idx := 0
	newData := make([][]int64, 0, 1440)
	for i := 0; i < 1440; i++ {
		if idx >= len(oldData) {
			break
		}
		data := oldData[idx]
		ts := fst.Unix() + int64(i) * 60
		if data[timeIdx] > ts {
			newData = append(newData, []int64{fst.Add(time.Duration(i) * time.Minute).Unix(), -1, -1, -1, -1})
		} else {
			if data[timeIdx] < ts {
				i--
			} else {
				newData = append(newData, data)
			}
			idx++
		}
	}
	for i := len(newData); i < 1440; i++ {
		newData = append(newData, []int64{fst.Add(time.Duration(i) * time.Minute).Unix(), -1, -1, -1, -1})
	}
	d.DayData = newData
}

func (d *DataSet) Put(record *common.Record) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Buffer = record
}

func (d *DataSet) Refresh(now time.Time) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	cur := time.Unix(0, now.UnixNano()-now.UnixNano()%int64(time.Minute))
	for i := 0; i < 1439; i++ {
		d.DayData[i] = d.DayData[i+1]
	}
	if d.Buffer != nil {
		d.DayData[1439] = []int64{cur.Unix(), d.Buffer.Avg, d.Buffer.Max, d.Buffer.Min, d.Buffer.Lost}
	} else {
		d.DayData[1439] = []int64{cur.Unix(), -1, -1, -1, -1}
	}
	d.Buffer = nil
}

func (d *DataSet) GetAbbrData() *AbbrDataSet {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	data := make([][]int64, len(d.DayData))
	for idx, record := range d.DayData {
		data[idx] = []int64{record[timeIdx], record[avgIdx]}
	}
	return &AbbrDataSet{
		Observer: d.Observer,
		Data:     data,
	}
}

func (d *DataSet) GetFullData() *FullDataSet {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	data := make([][]int64, len(d.DayData))
	for idx, record := range d.DayData {
		data[idx] = record
	}
	return &FullDataSet{
		Observer: d.Observer,
		Data:     data,
	}
}
