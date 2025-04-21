package main

import (
	stability "github.com/BeroKiTeer/KitBridge/kitex_gen/thrift/stability/stservice"
	"log"
)

func main() {
	svr := stability.NewServer(new(STServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
