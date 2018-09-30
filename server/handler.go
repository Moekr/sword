package server

import (
	"encoding/json"
	"fmt"
	"github.com/Moekr/sword/common"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

func httpConf(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}
	if bs, err := json.Marshal(conf.Targets); err != nil {
		http.Error(w, "marshal conf error: "+err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(bs)
	}
}

func httpData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "http method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !checkToken(w, r) {
		return
	}
	var targetId, observerId int64
	targetId, err := parseIntParam(r, "t", false, -1)
	observerId, err = parseIntParam(r, "o", false, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read post data error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	record := &common.Record{}
	if err := json.Unmarshal(bs, record); err != nil {
		http.Error(w, "unmarshal port data error: "+err.Error(), http.StatusBadRequest)
		return
	}
	if v := dataSets[targetId]; v != nil {
		dataSet := v[observerId]
		if dataSet != nil {
			dataSet.Put(record)
			return
		}
	}
	http.Error(w, "no such target id or observer id", http.StatusBadRequest)
}

func checkToken(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get(common.TokenHeaderName) != args.Token {
		http.Error(w, "token invalid", http.StatusForbidden)
		return false
	}
	return true
}

func httpAbbrData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse form error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	targetId, err := parseIntParam(r, "t", false, -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := dataSets[targetId]
	result := make([]*AbbrDataSet, 0, len(data))
	for _, dataSet := range data {
		result = append(result, dataSet.GetAbbrData())
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Observer.Id < result[j].Observer.Id
	})
	if bs, err := json.Marshal(result); err != nil {
		http.Error(w, "marshal result error: "+err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(bs)
	}
}

func httpFullData(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

var errMissingParam = fmt.Errorf("missing required param")

func parseIntParam(r *http.Request, key string, nullable bool, defaultValue int64) (int64, error) {
	r.ParseForm()
	str := r.Form.Get(key)
	if str == "" {
		if nullable {
			return defaultValue, nil
		} else {
			return defaultValue, errMissingParam
		}
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultValue, err
	}
	return val, nil
}
