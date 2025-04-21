package dal

import (
	"github.com/BeroKiTeer/KitBridge/biz/dal/mysql"
	"github.com/BeroKiTeer/KitBridge/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
