package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/huannet/IrisQrcode/controllers"
	"github.com/gomodule/redigo/redis"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"github.com/kataras/iris/middleware/basicauth"
	"time"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8090", "the port to listen for api service")
}

func main() {
	flag.Parse()
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := iris.New()
	app.Logger().SetLevel(os.Getenv("LOG_LEVEL"))
	app.Use(recover.New())
	app.Use(logger.New())

	app.Macros().String.RegisterFunc("urlCode", func(lent int) func(string) bool {
		return func(paramValue string) bool {
			return len(paramValue) == lent
		}
	})
	mvc.Configure(app.Party("/"), basicMVC)
	mvc.Configure(app.Party("/admin", basicauth.New(basicauth.Config{
		Users:   map[string]string{"guojiayi": "guojiayi#123", "zhanghao": "zhanghao#0807"},
		Realm:   "Authorization Required",
		Expires: time.Duration(30) * time.Minute,
	})), adminMVC)

	app.Run(iris.Addr(":8080"))
}

func basicMVC(app *mvc.Application) {
	app.Register(NewRedisPool())
	app.Handle(new(controllers.ApiController))
}

func adminMVC(app *mvc.Application) {
	app.Register(NewRedisPool())
	app.Handle(new(controllers.AdminController))
}

func NewRedisPool() *redis.Pool {
	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	poolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOLSIZE"))
	if poolSize < 100 {
		poolSize = 100
	}
	return &redis.Pool{
		MaxIdle:   20,
		MaxActive: poolSize,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", os.Getenv("REDIS_ADDR"))
			if err != nil {
				panic(err.Error())
				return nil, err
			}
			if os.Getenv("REDIS_PASSWD") != "" {
				if _, err := c.Do("AUTH", os.Getenv("REDIS_PASSWD")); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", redisDb); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}
}