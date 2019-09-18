package ping

import (
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/Moekr/gopkg/must"
	"github.com/Moekr/sword/types"
)

func Ping(address string) *types.TRecord {
	rec := types.NewRecord()
	conn, err := net.Dial("ip4:icmp", address)
	if err != nil {
		return rec
	}
	defer must.Close(conn)

	buf, vs := make([]byte, 1024), make([]float64, 0, 20)
	for seq := 0; seq < 20; seq++ {
		req := newRequest(uint16(seq))
		start := time.Now()
		if _, err := conn.Write(req); err != nil {
			continue
		} else if err := conn.SetReadDeadline(start.Add(2 * time.Second)); err != nil {
			continue
		} else if length, err := conn.Read(buf); err != nil {
			continue
		} else if length != 20+64 || !checkResponse(req, buf[20:84]) {
			continue
		}
		vs = append(vs, float64(time.Now().Sub(start).Nanoseconds())/float64(time.Millisecond))
		time.Sleep(250 * time.Millisecond)
	}

	return types.BuildRecord(vs, 20)
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

func checkResponse(req, rsp []byte) bool {
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
