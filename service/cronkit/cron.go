package cronkit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/commonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

//
var pool map[string]*cron.Cron

func Scheduler() *cron.Cron {
	if scheduler == nil {
		scheduler = NewScheduler()
	}
	return scheduler
}

func NewScheduler() *cron.Cron {
	return cron.New(cron.WithSeconds(), cron.WithLocation(configkit.GetLocation()))
}

func AddPool(key string, cron *cron.Cron) {
	RemovePool(key)
	pool[key] = cron
}

func RemovePool(key string) {
	v, ok := pool[key]
	if ok {
		v.Stop()
		delete(pool, key)
	}
}

// AddFunc 给默认的scheduler add func， 封装上recover
func AddFunc(spec string, fun func()) {
	_, err := Scheduler().AddFunc(spec, func() {
		commonkit.RecoverFuncWrapper(fun)
	})
	if err != nil {
		panic(exception.New(err.Error()))
	}
}
