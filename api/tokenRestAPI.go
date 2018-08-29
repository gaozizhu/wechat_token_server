package api

import(

	"wechat-tokenServer/tokenserver"
	"wechat-tokenServer/logger"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

var(
	tredisURL, tappid, tappSecret string
)

func StartRestServer(redisURL, appid, appSecret string)  {
	tredisURL = redisURL
	tappid = appid
	tappSecret = appSecret

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router,err := rest.MakeRouter(
		rest.Get("/accesstoken",GetAccessTokenFromRedis),
		rest.Post("/accesstoken/#token/#invalidtime",ModifyAccessTokenInRedis),
	)
	if err != nil {
		logger.Error.Println(err)
	}
	api.SetApp(router)
	http.ListenAndServe(":8777",api.MakeHandler())

	logger.Trace.Println("/accesstoken" + " ---> 从数据库中获取accesstoken")
	logger.Trace.Println("/accesstoken/:act" + " --->修改数据库中的accesstoken")
	logger.Trace.Println("API服务器已经启动，监听端口8777,服务名称/getaccesstoken")


}

/**
从Redis中获取AccessToken
 */
func GetAccessTokenFromRedis(w rest.ResponseWriter, r *rest.Request) {
	acessToken := tokenserver.ResAccessToken{}
	acessToken.GetAccessTokenFromRedis(tredisURL, tappid, tappSecret)
	w.WriteJson(acessToken)

}


//修改Redis中的accesstoken
func ModifyAccessTokenInRedis(w rest.ResponseWriter, r *rest.Request)  {
	accesstoken := r.PathParam("token")
	invalidtime := r.PathParam("invalidtime")
	acessToken := tokenserver.ResAccessToken{}
	acessToken.ModifyAccessTokenInRedis(tredisURL,tappid,tappSecret,accesstoken,invalidtime)
	w.WriteJson(&accesstoken)
}
