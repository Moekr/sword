package client

import (
	"github.com/Moekr/sword/common"
	"github.com/Moekr/sword/util/logs"
	"math"
	"math/rand"
	"net"
	"os"
	"time"
)

func doPing(address string) *common.Record {
	now := time.Now()
	cur := time.Unix(0, now.UnixNano()-now.UnixNano()%int64(time.Minute))
	record := &common.Record{
		Time: cur.Unix(),
		Avg:  -1,
		Max:  -1,
		Min:  -1,
		Lost: -1,
	}

	conn, err := net.Dial("ip4:icmp", address)
	if err != nil {
		logs.Debug("dial %s error: %s", address, err.Error())
		return record
	}
	defer conn.Close()

	var totalLatency, maxLatency, minLatency, lostPacket int64
	minLatency = math.MaxInt64
	var sequence uint16
	buffer := make([]byte, 1024)
	for sequence = 0; sequence < 20; sequence++ {
		req := newRequest(sequence)
		if _, err := conn.Write(req); err != nil {
			logs.Debug("write request to %s error: %s", address, err.Error())
			lostPacket++
			continue
		}
		start := time.Now()
		conn.SetReadDeadline(start.Add(2 * time.Second))
		length, err := conn.Read(buffer)
		if err != nil {
			logs.Debug("read response from %s error: %s", address, err.Error())
			lostPacket++
			continue
		} else if length != 20+64 || !validateResponse(req, buffer[20:84]) {
			logs.Debug("response from %s invalid", address)
			lostPacket++
			continue
		}
		end := time.Now()

		latency := end.Sub(start).Nanoseconds() / int64(time.Millisecond)
		logs.Debug("response from %s latency %dms", address, latency)
		totalLatency = totalLatency + latency
		if latency > maxLatency {
			maxLatency = latency
		}
		if latency < minLatency {
			minLatency = latency
		}
		time.Sleep(250 * time.Millisecond)
	}
	if lostPacket < 20 {
		record.Avg = totalLatency / (20 - lostPacket)
		record.Max = maxLatency
		record.Min = minLatency
	}
	record.Lost = lostPacket * 5
	return record
}

func newRequest(sequence uint16) []byte {
	req := make([]byte, 64)
	req[0] = 8
	pid := os.Getpid()
	req[4] = byte(pid >> 8)
	req[5] = byte(pid & 0xff)
	req[6] = byte(sequence >> 8)
	req[7] = byte(sequence & 0xff)
	for i := 8; i < 64; i++ {
		req[i] = byte(rand.Int() & 0xff)
	}
	cs := checkSum(req)
	req[2] = byte(cs >> 8)
	req[3] = byte(cs & 0xff)
	return req
}

func validateResponse(req, rsp []byte) bool {
	if rsp[0] != 0 || rsp[1] != 0 {
		return false
	}
	for i := 4; i < 64; i++ {
		if req[i] != rsp[i] {
			return false
		}
	}
	cs := (uint16(rsp[2]) << 8) + uint16(rsp[3])
	rsp[2] = 0
	rsp[3] = 0
	return checkSum(rsp) == cs
}

func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16
	return uint16(^sum)
}
