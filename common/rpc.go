package common

const (
	TokenHeaderName = "X-Sword-Token"
)

type Target struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Observer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Conf struct {
	Targets   []*Target   `json:"targets"`
	Observers []*Observer `json:"observers"`
}

type Record struct {
	Time int64 `json:"time"`
	Avg  int64 `json:"avg"`
	Max  int64 `json:"max"`
	Min  int64 `json:"min"`
	Lost int64 `json:"lost"`
}
