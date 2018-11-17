package server

import (
	"github.com/Moekr/sword/common"
	"sort"
)

type Category struct {
	Id   int64       `json:"id"`
	Name string      `json:"name"`
	Sub  []*Category `json:"sub"`
}

type Observer struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Conf struct {
	DefaultCid int64            `json:"default_cid"`
	Categories []*Category      `json:"categories"`
	Targets    []*common.Target `json:"targets"`
	Observers  []*Observer      `json:"observers"`
}

func (c *Conf) Init() {
	for _, category := range c.Categories {
		category.Init()
	}
	sort.Slice(c.Targets, func(i, j int) bool {
		return c.Targets[i].Name < c.Targets[j].Name
	})
	sort.Slice(c.Observers, func(i, j int) bool {
		return c.Observers[i].Name < c.Observers[j].Name
	})
}

func (c *Category) Init() {
	if c.Sub == nil {
		return
	}
	sort.Slice(c.Sub, func(i, j int) bool {
		return c.Sub[i].Name < c.Sub[j].Name
	})
}

func (c *Conf) GetCategory(cid int64) *Category {
	return getCategory(cid, c.Categories)
}

func getCategory(cid int64, categories []*Category) *Category {
	for _, category := range categories {
		if category.Id == cid {
			return category
		}
		if sub := getCategory(cid, category.Sub); sub != nil {
			return sub
		}
	}
	return nil
}

func (c *Conf) GetTargets(cid int64) []*common.Target {
	targets := make([]*common.Target, 0)
	for _, target := range c.Targets {
		for _, val := range target.Cid {
			if val == cid {
				targets = append(targets, target)
				break
			}
		}
	}
	return targets
}
