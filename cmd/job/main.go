package job

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/1024casts/snake/cmd/job/demo"
	"github.com/1024casts/snake/pkg/log"
)

// 计划任务
// see: https://mp.weixin.qq.com/s/Ak7RBv1NuS-VBeDNo8_fww
//
// cron 内置3个用得比较多的JobWrapper：
//
// Recover：捕获内部Job产生的 panic；
// DelayIfStillRunning：触发时，如果上一次任务还未执行完成（耗时太长），则等待上一次任务完成之后再执行；
// SkipIfStillRunning：触发时，如果上一次任务还未完成，则跳过此次执行。
func main() {
	c := cron.New()
	// demo
	_, err := c.AddFunc("* */5 * * *", func() {
		log.Infof("test cron, time: %d ", time.Now().Unix())
	})
	if err != nil {
		log.Warnf("cron AddFunc err, %+v", err)
		return
	}

	// test recover
	c.AddJob("@every 1s", cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(&demo.PanicJob{}))

	// test DelayIfStillRunning
	c.AddJob("@every 1s", cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).Then(&demo.DelayJob{}))

	// test SkipIfStillRunning
	c.AddJob("@every 1s", cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(&demo.SkipJob{}))

	// 执行具体的任务
	c.AddJob("@every 3s", demo.GreetingJob{"dj"})

	c.Start()
}
