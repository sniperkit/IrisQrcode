package controllers

import (
	"github.com/kataras/iris"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/mvc"
)

type AdminController struct {
	Ctx       iris.Context
	RedisPool *redis.Pool
}

func (c AdminController) Get() mvc.Result {
	return mvc.View{
		Name: "views/test.html",
		Data: nil,
	}
	// return map[string]string{"message": "Hello Iris!"}
}
