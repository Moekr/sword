package dataset

import (
	"bytes"
	"time"

	"github.com/Moekr/gopkg/algo"
	"github.com/Moekr/sword/types"
)

const (
	headerLength   = 16
	dsHeaderLength = 16
	dsDataLength   = 9 * 1440 * 4
	dsLength       = dsHeaderLength + dsDataLength
)

func Decode(bs []byte) {
	if len(bs) < headerLength {
		return
	}
	header, bs := bs[:headerLength], bs[headerLength:]
	ots, nts := decodeInt32(header[:4]), int32(time.Now().Unix())
	var data []byte
	for len(bs) >= dsLength {
		header, data, bs = bs[:dsHeaderLength], bs[dsHeaderLength:dsLength], bs[dsLength:]
		tid, cid := decodeInt16(header[:2]), decodeInt16(header[2:4])
		if tds, ok := DataSets[tid]; ok {
			if cds, ok := tds.Data[cid]; ok {
				decodeDataSet(data, cds.Data)
				shiftDataSet(ots, nts, cds.Data)
			}
		}
	}
}

func decodeDataSet(bs []byte, ds *types.TDataSet) {
	var rs []*types.TRecord
	rs = append(rs, ds.Day...)
	rs = append(rs, ds.Week...)
	rs = append(rs, ds.Month...)
	rs = append(rs, ds.Year...)
	var data []byte
	for _, record := range rs {
		data, bs = bs[:9], bs[9:]
		record.Avg = decodeInt16(data[:2])
		record.Max = decodeInt16(data[2:4])
		record.Min = decodeInt16(data[4:6])
		record.Std = decodeInt16(data[6:8])
		record.Los = int8(data[8])
	}
}

func shiftDataSet(ots, nts int32, ds *types.TDataSet) {
	shift := func(tick int32, rs []*types.TRecord) {
		offset := algo.MinInt(int(nts/tick-ots/tick), 1440)
		if offset > 0 {
			for i := 0; i < 1440-offset; i++ {
				rs[i].Avg = rs[i+offset].Avg
				rs[i].Max = rs[i+offset].Max
				rs[i].Min = rs[i+offset].Min
				rs[i].Std = rs[i+offset].Std
				rs[i].Los = rs[i+offset].Los
			}
			for i := 1440 - offset; i < 1440; i++ {
				rs[i].Avg = -1
				rs[i].Max = -1
				rs[i].Min = -1
				rs[i].Std = -1
				rs[i].Los = -1
			}
		}
	}
	shift(60, ds.Day)
	shift(7*60, ds.Week)
	shift(30*60, ds.Month)
	shift(360*60, ds.Year)
}

func Encode() []byte {
	buf := &bytes.Buffer{}
	// header
	buf.Write(encodeInt32(int32(time.Now().Unix())))
	emptyBytes(buf, 12)
	// ds
	for tid, tds := range DataSets {
		for cid, cds := range tds.Data {
			buf.Grow(dsLength)
			// dsHeader
			buf.Write(encodeInt16(tid))
			buf.Write(encodeInt16(cid))
			emptyBytes(buf, 12)
			// dsData
			encodeDataSet(buf, cds.Data)
		}
	}
	return buf.Bytes()
}

func encodeDataSet(buf *bytes.Buffer, ds *types.TDataSet) {
	var rs []*types.TRecord
	rs = append(rs, ds.Day...)
	rs = append(rs, ds.Week...)
	rs = append(rs, ds.Month...)
	rs = append(rs, ds.Year...)
	for _, record := range rs {
		buf.Write(encodeInt16(record.Avg))
		buf.Write(encodeInt16(record.Max))
		buf.Write(encodeInt16(record.Min))
		buf.Write(encodeInt16(record.Std))
		buf.WriteByte(byte(record.Los))
	}
}

func decodeInt16(bs []byte) int16 {
	return int16(bs[0])<<8 + int16(bs[1])
}

func decodeInt32(bs []byte) int32 {
	return int32(bs[0])<<24 + int32(bs[1])<<16 + int32(bs[2])<<8 + int32(bs[3])
}

func decodeInt64(bs []byte) int64 {
	return int64(bs[0])<<56 + int64(bs[1])<<48 + int64(bs[2])<<40 + int64(bs[3])<<32 +
		int64(bs[4])<<24 + int64(bs[5])<<16 + int64(bs[6])<<8 + int64(bs[7])
}

func encodeInt16(i16 int16) []byte {
	return []byte{
		byte((i16 >> 8) & 0xff), byte((i16 >> 0) & 0xff),
	}
}

func encodeInt32(i32 int32) []byte {
	return []byte{
		byte((i32 >> 24) & 0xff), byte((i32 >> 16) & 0xff), byte((i32 >> 8) & 0xff), byte((i32 >> 0) & 0xff),
	}
}

func encodeInt64(i64 int64) []byte {
	return []byte{
		byte((i64 >> 56) & 0xff), byte((i64 >> 48) & 0xff), byte((i64 >> 40) & 0xff), byte((i64 >> 32) & 0xff),
		byte((i64 >> 24) & 0xff), byte((i64 >> 16) & 0xff), byte((i64 >> 8) & 0xff), byte((i64 >> 0) & 0xff),
	}
}

func emptyBytes(buf *bytes.Buffer, c int) {
	for i := 0; i < c; i++ {
		buf.WriteByte(0)
	}
}
