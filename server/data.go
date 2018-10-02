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

const (
	_ = iota
	rangeDay
	rangeWeek
	rangeMonth
	rangeYear
)

const (
	dayInterval   = time.Minute
	weekInterval  = 7 * time.Minute
	monthInterval = 30 * time.Minute
	yearInterval  = 6 * time.Hour
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
	d.DayData = initData(d.DayData, dayInterval, 1440)
	d.WeekData = initData(d.WeekData, weekInterval, 1440)
	d.MonthData = initData(d.MonthData, monthInterval, 1440)
	d.YearData = initData(d.YearData, yearInterval, 1440)
}

func initData(origin [][]int64, interval time.Duration, count int) [][]int64 {
	cur := time.Unix(0, initNow.UnixNano()-initNow.UnixNano()%int64(interval))
	fst := cur.Add(-time.Duration(count) * interval)
	if interval > dayInterval {
		fst = fst.Add(interval)
	}
	oldData := make([][]int64, 0)
	for idx, data := range origin {
		if data[timeIdx] == fst.Unix() {
			oldData = origin[idx:]
			break
		}
	}
	idx := 0
	newData := make([][]int64, 0, count)
	for i := 0; i < count; i++ {
		if idx >= len(oldData) {
			break
		}
		data := oldData[idx]
		ts := fst.Unix() + int64(i)*int64(interval.Seconds())
		if data[timeIdx] > ts {
			newData = append(newData, []int64{fst.Add(time.Duration(i) * interval).Unix(), -1, -1, -1, -1})
		} else {
			if data[timeIdx] < ts {
				i--
			} else {
				newData = append(newData, data)
			}
			idx++
		}
	}
	for i := len(newData); i < count; i++ {
		newData = append(newData, []int64{fst.Add(time.Duration(i) * interval).Unix(), -1, -1, -1, -1})
	}
	return newData
}

func (d *DataSet) Put(record *common.Record) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Buffer = record
}

func (d *DataSet) Refresh(cur time.Time) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	for i := 0; i < 1439; i++ {
		d.DayData[i] = d.DayData[i+1]
	}
	if d.Buffer != nil {
		d.DayData[1439] = []int64{cur.Unix(), d.Buffer.Avg, d.Buffer.Max, d.Buffer.Min, d.Buffer.Lost}
	} else {
		d.DayData[1439] = []int64{cur.Unix(), -1, -1, -1, -1}
	}
	if cur.UnixNano()%int64(weekInterval) == 0 {
		for i := 0; i < 1439; i++ {
			d.WeekData[i] = d.WeekData[i+1]
		}
		d.WeekData[1439] = average(d.DayData[1433:])
	}
	if cur.UnixNano()%int64(monthInterval) == 0 {
		for i := 0; i < 1439; i++ {
			d.MonthData[i] = d.MonthData[i+1]
		}
		d.MonthData[1439] = average(d.DayData[1410:])
	}
	if cur.UnixNano()%int64(yearInterval) == 0 {
		for i := 0; i < 1439; i++ {
			d.YearData[i] = d.YearData[i+1]
		}
		d.YearData[1439] = average(d.DayData[1080:])
	}
	d.Buffer = nil
}

func average(data [][]int64) []int64 {
	result := make([]int64, length)
	result[timeIdx] = data[len(data)-1][timeIdx]
	var avg, max, min, lost, cnt int64 = -1, -1, -1, -1, 0
	for _, d := range data {
		if d[avgIdx] == -1 || d[maxIdx] == -1 || d[minIdx] == -1 || d[lostIdx] == -1 {
			continue
		}
		avg = avg + d[avgIdx]
		max = max + d[maxIdx]
		min = min + d[minIdx]
		lost = lost + d[lostIdx]
		cnt++
	}
	if cnt == 0 {
		cnt = 1
	}
	result[avgIdx] = avg / cnt
	result[maxIdx] = max / cnt
	result[minIdx] = min / cnt
	result[lostIdx] = lost / cnt
	return result
}

func (d *DataSet) GetAbbrData(timeRange int64) *AbbrDataSet {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	originData := d.GetOriginData(timeRange)
	data := make([][]int64, len(originData))
	for idx, record := range originData {
		data[idx] = []int64{record[timeIdx], record[avgIdx]}
	}
	return &AbbrDataSet{
		Observer: d.Observer,
		Data:     data,
	}
}

func (d *DataSet) GetFullData(timeRange int64) *FullDataSet {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	originData := d.GetOriginData(timeRange)
	data := make([][]int64, len(originData))
	for idx, record := range originData {
		data[idx] = record
	}
	return &FullDataSet{
		Observer: d.Observer,
		Data:     data,
	}
}

func (d *DataSet) GetOriginData(timeRange int64) [][]int64 {
	switch timeRange {
	case rangeDay:
		return d.DayData
	case rangeWeek:
		return d.WeekData
	case rangeMonth:
		return d.MonthData
	case rangeYear:
		return d.YearData
	}
	return d.DayData
}
