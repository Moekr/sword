package server

import "github.com/Moekr/sword/common"

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
