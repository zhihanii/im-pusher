package broker

import (
	"github.com/zhihanii/im-pusher/internal/broker/conf"
	"os"
	"runtime/pprof"
)

func Run(cfg *conf.Config) error {
	s, err := createServer(cfg)
	if err != nil {
		return err
	}
	cpuProfile, err := os.OpenFile("/opt/broker/pprof/broker-cpu.profile", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cpuProfile.Close()
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	memProfile, err := os.OpenFile("/opt/broker/pprof/broker-mem.profile", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer memProfile.Close()
	defer pprof.WriteHeapProfile(memProfile)

	return s.Prepare().Run()
}
