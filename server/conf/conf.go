package conf

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Moekr/sword/common/args"
	"github.com/Moekr/sword/types"
)

var (
	Conf       *types.TConf
	Categories map[int16]*types.TCategory
	Targets    map[int16]*types.TTarget
	C2Targets  map[int16][]*types.TTarget
)

func InitConf() error {
	if err := loadConf(); err != nil {
		return err
	}
	Categories = make(map[int16]*types.TCategory)
	for _, category := range Conf.Categories {
		initCategory(category)
	}
	sort.Slice(Conf.Targets, func(i, j int) bool {
		return Conf.Targets[i].Name < Conf.Targets[j].Name
	})
	Targets, C2Targets = make(map[int16]*types.TTarget), make(map[int16][]*types.TTarget)
	for _, target := range Conf.Targets {
		Targets[target.ID] = target
		for _, cid := range target.Cid {
			C2Targets[cid] = append(C2Targets[cid], target)
		}
	}
	sort.Slice(Conf.Clients, func(i, j int) bool {
		return Conf.Clients[i].Name < Conf.Clients[j].Name
	})
	return nil
}

func loadConf() error {
	confPath := args.Args.ConfPath
	if bs, err := ioutil.ReadFile(confPath); err != nil {
		return err
	} else {
		confPath = strings.ToLower(confPath)
		if strings.HasSuffix(confPath, ".yml") || strings.HasSuffix(confPath, ".yaml") {
			return yaml.Unmarshal(bs, &Conf)
		}
		return json.Unmarshal(bs, &Conf)
	}
}

func initCategory(category *types.TCategory) {
	Categories[category.ID] = category
	if category.Sub == nil {
		return
	}
	for _, category := range category.Sub {
		Categories[category.ID] = category
	}
	sort.Slice(category.Sub, func(i, j int) bool {
		return category.Sub[i].Name < category.Sub[j].Name
	})
}
