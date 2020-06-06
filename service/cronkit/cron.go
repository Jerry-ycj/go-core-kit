package cronkit

import (
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

//
var pool map[string]*cron.Cron

func Scheduler() *cron.Cron {
	if scheduler == nil {
		scheduler = New()
	}
	return scheduler
}

func New() *cron.Cron {
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

/**
note: https://godoc.org/github.com/robfig/cron

c := cron.New()
c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
c.AddFunc("@hourly",      func() { fmt.Println("Every hour, starting an hour from now") })
c.AddFunc("@every 1h30m10s", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
c.Start()
..
// Funcs are invoked in their own goroutine, asynchronously.
...
// Funcs may also be added to a running Cron
c.AddFunc("@daily", func() { fmt.Println("Every day") })
..
// Inspect the cron job entries' next and previous run times.
inspect(c.Entries())
..
c.Stop()  // Stop the scheduler (does not stop any jobs already running).

// cron format
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

// job wrappers
- Recover any panics from jobs (activated by default)
- Delay a job's execution if the previous run hasn't completed yet
- Skip a job's execution if the previous run hasn't completed yet
- Log each job's invocations
*/
