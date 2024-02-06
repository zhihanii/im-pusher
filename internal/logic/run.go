package logic

import (
	"github.com/zhihanii/im-pusher/internal/logic/conf"
	"os"
	"runtime/pprof"
)

func Run(c *conf.Config) error {
	s, err := createServer(c)
	if err != nil {
		return err
	}
	cpuProfile, err := os.OpenFile("/opt/logic/pprof/logic-cpu.profile", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cpuProfile.Close()
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	memProfile, err := os.OpenFile("/opt/logic/pprof/logic-mem.profile", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer memProfile.Close()
	defer pprof.WriteHeapProfile(memProfile)

	return s.Prepare().Run()
}
