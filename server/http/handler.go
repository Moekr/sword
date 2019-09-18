package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"gopkg.in/macaron.v1"

	"github.com/Moekr/sword/common/version"
	"github.com/Moekr/sword/server/conf"
	"github.com/Moekr/sword/server/dataset"
	"github.com/Moekr/sword/types"
)

func handleIndexPage(ctx *macaron.Context) {
	cid, timeRange := int16(ctx.QueryInt("c")), formatTimeRange(ctx.QueryInt("r"))
	if ctx.Query("c") == "" {
		cid = conf.Conf.Index
	}
	category := conf.Categories[cid]
	if category == nil {
		ctx.Redirect("/index.html", http.StatusMovedPermanently)
	} else {
		ctx.HTML(http.StatusOK, "index", map[string]interface{}{
			"categories": conf.Conf.Categories,
			"category":   category,
			"targets":    conf.C2Targets[cid],
			"timeRange":  timeRange,
			"version":    version.Version,
		})
	}
}

func handleDetailPage(ctx *macaron.Context) {
	tid, timeRange := int16(ctx.QueryInt("t")), formatTimeRange(ctx.QueryInt("r"))
	target := conf.Targets[tid]
	if target == nil {
		ctx.Redirect("/index.html", http.StatusMovedPermanently)
	} else {
		ctx.HTML(http.StatusOK, "detail", map[string]interface{}{
			"categories": conf.Conf.Categories,
			"target":     target,
			"clients":    conf.Conf.Clients,
			"timeRange":  timeRange,
			"version":    version.Version,
		})
	}
}

func handleAPIConf() ([]byte, error) {
	return json.Marshal(newDataResponse(conf.Conf.Targets))
}

func handleAPIPush(ctx *macaron.Context) {
	tid, cid := int16(ctx.QueryInt("t")), int16(ctx.QueryInt("c"))
	var ds *types.TDataSet
	if tds, ok := dataset.DataSets[tid]; ok {
		if cds, ok := tds.Data[cid]; ok {
			ds = cds.Data
		}
	}
	if ds == nil {
		ctx.Error(http.StatusNotFound)
	} else if bs, err := ioutil.ReadAll(ctx.Req.Request.Body); err != nil {
		ctx.Error(http.StatusInternalServerError)
	} else {
		var record *types.TRecord
		if err := json.Unmarshal(bs, &record); err != nil {
			ctx.Error(http.StatusBadRequest)
		} else {
			ds.Lock.Lock()
			defer ds.Lock.Unlock()
			ds.Buf = record
		}
	}
}

func handleAPIAbbrQuery(ctx *macaron.Context) {
	tid, timeRange := int16(ctx.QueryInt("t")), formatTimeRange(ctx.QueryInt("r"))
	tds, ok := dataset.DataSets[tid]
	if !ok {
		ctx.Error(http.StatusNotFound)
		return
	}
	res := make([]*QueryData, 0, len(tds.Data))
	for _, cds := range tds.Data {
		res = append(res, newAbbrQueryData(cds, timeRange))
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Client.Name < res[j].Client.Name
	})
	encoder := json.NewEncoder(ctx.Resp)
	if err := encoder.Encode(newDataResponse(res)); err != nil {
		ctx.Error(http.StatusInternalServerError)
	}
}

func handleAPIFullQuery(ctx *macaron.Context) {
	tid, cid, timeRange := int16(ctx.QueryInt("t")), int16(ctx.QueryInt("c")), formatTimeRange(ctx.QueryInt("r"))
	tds, ok := dataset.DataSets[tid]
	if !ok {
		ctx.Error(http.StatusNotFound)
		return
	}
	cds, ok := tds.Data[cid]
	if !ok {
		ctx.Error(http.StatusNotFound)
		return
	}
	res := newFullQueryData(cds, timeRange)
	encoder := json.NewEncoder(ctx.Resp)
	if err := encoder.Encode(newDataResponse(res)); err != nil {
		ctx.Error(http.StatusInternalServerError)
	}
}

func handleAPIStat(ctx *macaron.Context) {
	tid, interval := int16(ctx.QueryInt("t")), ctx.QueryInt("i")
	tds, ok := dataset.DataSets[tid]
	if !ok {
		ctx.Error(http.StatusNotFound)
		return
	}
	res := &StatsData{
		Target: tds.Target,
		Stats:  make([]*StatData, 0, len(tds.Data)),
	}
	for _, cds := range tds.Data {
		res.Stats = append(res.Stats, newStatData(cds, interval))
	}
	sort.Slice(res.Stats, func(i, j int) bool {
		return res.Stats[i].Client.Name < res.Stats[j].Client.Name
	})
	encoder := json.NewEncoder(ctx.Resp)
	if err := encoder.Encode(newDataResponse(res)); err != nil {
		ctx.Error(http.StatusInternalServerError)
	}
}

func formatTimeRange(timeRange int) int {
	if timeRange < rangeDay || timeRange > rangeYear {
		return rangeDay
	}
	return timeRange
}
