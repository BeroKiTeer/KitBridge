package main

import (
	"github.com/BeroKiTeer/KitBridge/autodetect"
	"github.com/BeroKiTeer/KitBridge/http1"
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
	httpHandlerFactory := &http1.HTTP1SvrTransHandlerFactory{}
	opts = append(opts,
		server.WithTransHandlerFactory(autodetect.NewSvrTransHandlerFactoryWithHTTP(httpHandlerFactory)),
	)
	return
}
