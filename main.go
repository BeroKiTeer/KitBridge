package main

import (
	stability "github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability/stservice"
	"github.com/cloudwego/kitex/server"
	"log"
)

func main() {
	opts := kitexInit()

	svr := stability.NewServer(new(STServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func kitexInit() (opts []server.Option) {

	return
}
