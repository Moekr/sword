package http

import (
	"net/http"

	"github.com/Moekr/gopkg/algo"
	"github.com/Moekr/sword/server/dataset"
	"github.com/Moekr/sword/types"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type QueryData struct {
	Client *types.TClient `json:"client"`
	Avg    []int16        `json:"avg"`
	Max    []int16        `json:"max"`
	Min    []int16        `json:"min"`
	Std    []int16        `json:"std"`
	Los    []int8         `json:"los"`
}

type StatsData struct {
	Target *types.TTarget `json:"target"`
	Stats  []*StatData    `json:"stats"`
}

type StatData struct {
	Client *types.TClient `json:"client"`
	Avg    int16          `json:"avg"`
	Max    int16          `json:"max"`
	Min    int16          `json:"min"`
	Std    int16          `json:"std"`
	Los    int8           `json:"los"`
}

const (
	_ = iota
	rangeDay
	rangeWeek
	rangeMonth
	rangeYear
)

func newDataResponse(data interface{}) *Response {
	return &Response{
		Code: http.StatusOK,
		Data: data,
	}
}

func newErrorResponse(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

func newAbbrQueryData(cds *types.TCDataSet, timeRange int) *QueryData {
	data := &QueryData{
		Client: cds.Client,
		Avg:    make([]int16, 1440),
	}
	rs := selectRecords(cds, timeRange)
	cds.Data.Lock.RLock()
	defer cds.Data.Lock.RUnlock()
	for idx, rec := range rs {
		data.Avg[idx] = rec.Avg
	}
	return data
}

func newFullQueryData(cds *types.TCDataSet, timeRange int) *QueryData {
	data := &QueryData{
		Client: cds.Client,
		Avg:    make([]int16, 1440),
		Max:    make([]int16, 1440),
		Min:    make([]int16, 1440),
		Std:    make([]int16, 1440),
		Los:    make([]int8, 1440),
	}
	rs := selectRecords(cds, timeRange)
	cds.Data.Lock.RLock()
	defer cds.Data.Lock.RUnlock()
	for idx, rec := range rs {
		data.Avg[idx] = rec.Avg
		data.Max[idx] = rec.Max
		data.Min[idx] = rec.Min
		data.Std[idx] = rec.Std
		data.Los[idx] = rec.Los
	}
	return data
}

func newStatData(cds *types.TCDataSet, interval int) *StatData {
	data := &StatData{
		Client: cds.Client,
	}
	interval = algo.MaxInt(algo.MinInt(interval, 1440), 10)
	rec := dataset.Union(selectRecords(cds, rangeDay)[1440-interval:])
	data.Avg, data.Max, data.Min, data.Std, data.Los = rec.Avg, rec.Max, rec.Min, rec.Std, rec.Los
	return data
}

func selectRecords(cds *types.TCDataSet, timeRange int) []*types.TRecord {
	switch timeRange {
	case rangeDay:
		return cds.Data.Day
	case rangeWeek:
		return cds.Data.Week
	case rangeMonth:
		return cds.Data.Month
	case rangeYear:
		return cds.Data.Year
	}
	return cds.Data.Day
}
