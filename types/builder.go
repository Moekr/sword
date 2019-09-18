package types

import (
	"math"

	"github.com/Moekr/gopkg/algo"
)

func NewRecord() *TRecord {
	return &TRecord{
		Avg: -1,
		Max: -1,
		Min: -1,
		Std: -1,
		Los: -1,
	}
}

func BuildRecord(vs []float64, cnt int) *TRecord {
	rec := NewRecord()
	if len(vs) == 0 {
		rec.Los = 100
		return rec
	}
	sum, max, min := 0.0, 0.0, math.MaxFloat64
	for _, v := range vs {
		sum, max, min = sum+v, algo.MaxF64(max, v), algo.MinF64(min, v)
	}
	rec.Avg = int16(algo.Round(sum/float64(len(vs)), 0))
	rec.Max = int16(algo.Round(max, 0))
	rec.Min = int16(algo.Round(min, 0))
	rec.Std = int16(algo.Round(algo.Std(vs), 0))
	rec.Los = int8(100 * (cnt - len(vs)) / cnt)
	return rec
}
