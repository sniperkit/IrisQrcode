package controllers

import (
	"github.com/kataras/iris"
	"github.com/gomodule/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/huannet/IrisQrcode/models"
	"time"
	"os"
	"github.com/huannet/IrisQrcode/tools"
	"strings"
	"log"
	"github.com/kataras/iris/mvc"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
)

type ApiController struct {
	Ctx iris.Context
	RedisPool    *redis.Pool
}

func (c ApiController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/a/{code:string urlCode(8)}", "ShowUrlCode")
	b.Handle("GET", "/q/{code:string urlCode(8)}", "QrCodeByCode")
}

func (c ApiController) QrCodeByCode(){
	code := c.Ctx.Params().Get("code")
	qrValue := os.Getenv("HOME_URL")+"/a/"+code
	qrCode, _ := qr.Encode(qrValue, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 200, 200)
	filename := code+".png"
	fullname := "./storage/upload/"+filename
	exists , _ := tools.PathExists(fullname)
	if !exists {
		file, _ := os.Create(fullname)
		defer file.Close()
		png.Encode(file, qrCode)
	}
	c.Ctx.SendFile(fullname, filename)
}

func (c ApiController) ShowUrlCode() string {
	code := c.Ctx.Params().Get("code")

	redisClient := c.RedisPool.Get()
	defer redisClient.Close()

	cacheKey := "UrlCode_"+code
	cacheVal, _ := redis.String(redisClient.Do("GET", cacheKey))
	if cacheVal != "" {
		return cacheVal
	} else {
		return "url error"
	}
}

func (c ApiController) Get() interface{} {
	return iris.Map{
		"success": true,
		"msg":     "This is api page.",
	}
}

func (c ApiController) PostUrl() interface{} {
	urlStr := strings.TrimSpace(c.Ctx.PostValue("url"))
	if urlStr == "" {
		c.Ctx.StatusCode(500)
		return iris.Map{
			"success": false,
			"msg":     "url is null",
		}
	}

	db, dberr := gorm.Open("mysql", os.Getenv("MYSQL_CON"))
	if dberr != nil {
		c.Ctx.StatusCode(500)
		return iris.Map{
			"success": false,
			"msg":     dberr.Error(),
		}
	}
	defer db.Close()

	urlObj := models.Url{}
	dbResult := db.Where("to_url = ? ", urlStr).First(&urlObj)
	if dbResult.RecordNotFound() {
		urlCode := GetUrlCode(db, 0)
		if urlCode == "" {
			c.Ctx.StatusCode(500)
			return iris.Map{
				"success": false,
				"msg":     "url id generate error, please retry!",
			}
		}
		urlObj.Code = urlCode
		urlObj.ToUrl = urlStr
		dbResult := db.Create(&urlObj)
		if dbResult.Error != nil {
			c.Ctx.StatusCode(500)
			return iris.Map{
				"success": false,
				"msg":     "db NewRecord error",
			}
		} else {
			log.Println("new url record:", urlCode)
			c.CacheUrlObj(urlObj)
			return iris.Map{
				"success": true,
				"msg":     "ok",
			}
		}
	} else {
		c.Ctx.StatusCode(500)
		return iris.Map{
			"success": false,
			"msg":     "url is exists!",
		}
	}
}

func (c ApiController) GetTest() interface{} {

	redisClient := c.RedisPool.Get()
	defer redisClient.Close()

	cacheKey := "test"
	cacheVal := ""
	if exists, _:= redis.Bool(redisClient.Do("EXISTS", cacheKey)); exists {
		cacheVal, _ = redis.String(redisClient.Do("GET", cacheKey))
	} else {
		cacheVal = time.Now().Format("2006-01-02 15:04:05")
		redisClient.Do("SET", cacheKey, cacheVal, "EX", 10)
	}
	return iris.Map{
		"success": true,
		"msg":     cacheVal,
		"result":  "[]",
	}
}

func GetUrlCode(db *gorm.DB, count int) string {
	urlObj := models.Url{}
	urlCode := tools.GetRandomString(8);
	dbResult := db.Where("code = ? ", urlCode).First(&urlObj)
	if dbResult.RecordNotFound() {
		return urlCode
	} else {
		if count >= 3 {
			return ""
		} else {
			return GetUrlCode(db, count+1)
		}
	}
}

func (c ApiController) CacheUrlObj(urlObj models.Url) error {
	redisClient := c.RedisPool.Get()
	defer redisClient.Close()
	_, err := redisClient.Do("SET", "UrlCode_"+urlObj.Code, urlObj.ToUrl)
	return err
}

