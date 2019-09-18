package cronjob

import (
	"time"

	"github.com/Moekr/gopkg/logs"
	"github.com/Moekr/gopkg/periodic"
	"github.com/Moekr/sword/server/dataset"
	"github.com/Moekr/sword/server/persistence"
)

func StartCronJob() {
	periodic.NewStaticPeriodic(doCronJob, time.Minute, periodic.MinInterval).Start()
}

func doCronJob() {
	dataset.UpdateDataSets()
	if err := persistence.StoreData(true); err != nil {
		logs.Error("[CronJob] store data error: %s", err.Error())
	}
}
