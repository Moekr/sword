package client

import (
	"github.com/Moekr/sword/common"
	"log"
	"math"
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
		log.Println(err.Error())
		return record
	}
	defer conn.Close()

	pid := os.Getpid()
	req := make([]byte, 64)
	req[0] = 8
	req[4] = byte(pid >> 8)
	req[5] = byte(pid & 0xff)
	rsp := make([]byte, 1024)
	var totalLatency, maxLatency, minLatency, lostPacket int64
	minLatency = math.MaxInt64
	for i := 0; i < 20; i++ {
		req[2] = 0
		req[3] = 0
		req[7] = byte(i)
		cs := checkSum(req)
		req[2] = byte(cs >> 8)
		req[3] = byte(cs & 0xff)
		if _, err := conn.Write(req); err != nil {
			log.Println(err.Error())
			return record
		}
		start := time.Now()
		conn.SetReadDeadline(start.Add(5 * time.Second))
		_, err := conn.Read(rsp)
		if err != nil {
			log.Println(err.Error())
			lostPacket++
			continue
		}
		end := time.Now()
		latency := end.Sub(start).Nanoseconds() / int64(time.Millisecond)
		totalLatency = totalLatency + latency
		if latency > maxLatency {
			maxLatency = latency
		}
		if latency < minLatency {
			minLatency = latency
		}
		time.Sleep(300 * time.Millisecond)
	}
	if lostPacket < 20 {
		record.Avg = totalLatency / (20 - lostPacket)
		record.Max = maxLatency
		record.Min = minLatency
	}
	record.Lost = lostPacket * 5
	return record
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
