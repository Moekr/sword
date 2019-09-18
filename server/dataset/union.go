package dataset

import (
	"math"

	"github.com/Moekr/gopkg/algo"
	"github.com/Moekr/sword/types"
)

func Union(rs []*types.TRecord) *types.TRecord {
	res := types.NewRecord()
	var sum, max, min, stdSum, losSum, cnt int64 = 0, 0, math.MaxInt64, 0, 0, 0
	for _, rec := range rs {
		if rec.Los == -1 {
			continue
		}
		losSum = losSum + int64(rec.Los)
		if rec.Los == 100 {
			continue
		}
		c := (100 - int64(rec.Los)) / 5
		cnt = cnt + c
		sum = sum + int64(rec.Avg)*c
		max = algo.MaxI64(max, int64(rec.Max))
		min = algo.MinI64(min, int64(rec.Min))
		stdSum = stdSum + c*int64(rec.Avg)*int64(rec.Avg) + int64(rec.Std)*int64(rec.Std)
	}
	if cnt == 0 {
		sum, max, min, cnt = -1, -1, -1, 1
	}
	if cnt == 0 && losSum == 0 {
		res.Los = -1
	} else {
		res.Los = int8(losSum / int64(len(rs)))
	}
	avg := float64(sum) / float64(cnt)
	res.Avg = int16(algo.Round(avg, 0))
	res.Max = int16(max)
	res.Min = int16(min)
	res.Std = int16(algo.Round(math.Sqrt((float64(stdSum)-float64(cnt)*avg*avg)/float64(cnt)), 0))
	return res
}
