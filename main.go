package main

import (
	"wechat-tokenServer/api"
	"wechat-tokenServer/tokenserver"
)

const (
	//测试号
	//AppID     = "wx10a300ee87ff239f"
	//AppSecret = "845c55f125af95ae5045b974bc3ed615"

	//Redis地址和端口号
	RedisURL = "127.0.0.1:6379"

	//jsAuth跳转页面
	Redirect_uri = "http://192.168.8.100:8777/getjstokenCode"

	//万微公众号
	AppID     = "wxff8603fef038906d"
	AppSecret = "60cd515b54c64725737f73f8334639d4"
)

func main() {
	getToken()
}

func getToken() {

	//启动api服务器
	go api.StartRestServer(RedisURL, AppID, AppSecret)

	//启动jstoken服务器
	//go tokenserver.StartJSTokenServer(Redirect_uri, AppID, AppSecret)

	//从微信服务器获取Token
	//存入Redis中
	acessToken := tokenserver.ResAccessToken{}
	acessToken.InitResAccessToken(AppID, AppSecret) //从服务器初始化获得AccessTocken
	acessToken.InitRedisDB(RedisURL, AppID, AppSecret)

}
