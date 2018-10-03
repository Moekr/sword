package common

const (
	TokenHeaderName = "X-Sword-Token"
)

type Target struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Cid     []int64 `json:"cid"`
}

type Record struct {
	Time int64 `json:"time"`
	Avg  int64 `json:"avg"`
	Max  int64 `json:"max"`
	Min  int64 `json:"min"`
	Lost int64 `json:"lost"`
}
