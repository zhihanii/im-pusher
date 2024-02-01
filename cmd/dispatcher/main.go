package main

import (
	"github.com/zhihanii/im-pusher/internal/dispatcher"
)

func main() {
	dispatcher.NewApp("dispatcher").Run()
}
