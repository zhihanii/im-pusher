package main

import "github.com/zhihanii/im-pusher/internal/broker"

func main() {
	broker.NewApp("broker").Run()
}
