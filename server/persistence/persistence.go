package persistence

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Moekr/gopkg/compress"
	"github.com/Moekr/gopkg/fs"
	"github.com/Moekr/gopkg/logs"
	"github.com/Moekr/sword/common/args"
	"github.com/Moekr/sword/server/dataset"
)

func LoadData() ([]byte, error) {
	bs, err := ioutil.ReadFile(path.Join(args.Args.DataPath, "data.dat"))
	if err != nil {
		if !os.IsNotExist(err) {
			logs.Error("[Persistence] read data error: %s", err.Error())
		}
	} else if bs, err = compress.ZlibDecompress(bs); err != nil {
		logs.Error("[Persistence] decompress data error: %s", err.Error())
	}
	return bs, err
}

func StoreData(isCronJob bool) error {
	bs, err := compress.ZlibCompress(dataset.Encode())
	if err != nil {
		logs.Error("[Persistence] compress data error: %s", err.Error())
	} else if err = fs.AtomicWrite(path.Join(args.Args.DataPath, "data.dat"), bs, 0755); err != nil {
		logs.Error("[Persistence] write data error: %s", err.Error())
	} else if ts := time.Now().Unix() / 60 * 60; isCronJob && ts%3600 == 0 {
		bakPath := path.Join(args.Args.DataPath, fmt.Sprintf("data-%d.dat", ts))
		if err = fs.AtomicWrite(bakPath, bs, 0755); err != nil {
			logs.Error("[Persistence] write backup data error: %s", err.Error())
		}
	}
	return err
}
