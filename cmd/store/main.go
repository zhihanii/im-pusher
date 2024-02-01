package main

import "github.com/zhihanii/im-pusher/internal/store"

func main() {
	store.NewApp("store-server").Run()
}
