package api

import(

	"wechat-tokenServer/tokenserver"
	"wechat-tokenServer/logger"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"log"
	"net/url"
	"fmt"
)

var(
	tredisURL, tappid, tappSecret string
)
const(
	redirectOauthURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"

	//jsAuth跳转页面
	Redirect_uri = "http://gr.wanway.xin/getjstokenCode"
	)

//启动rest服务器
func StartRestServer(redisURL, appid, appSecret string)  {
	tredisURL = redisURL
	tappid = appid
	tappSecret = appSecret

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router,err := rest.MakeRouter(
		rest.Get("/accesstoken",GetAccessTokenFromRedis),
		rest.Get("/getjstokenurl",GetRedirectURL),
		rest.Get("/getjstokenCode",GetCodeforJStoken),
		rest.Post("/accesstoken/#token/#invalidtime",ModifyAccessTokenInRedis),
	)
	if err != nil {
		logger.Error.Println(err)
	}
	api.SetApp(router)
	log.Println("API服务器已经启动，监听端口8777,服务名称/getaccesstoken")
	http.ListenAndServe(":8777",api.MakeHandler())

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



/**
	获取code 步骤1-组装授权页URL
	appid				是	公众号的唯一标识
	redirect_uri		是	授权后重定向的回调链接地址，请使用urlencode对链接进行处理
	response_type		是	返回类型，请填写code
	scope				是	应用授权作用域，
							snsapi_base （不弹出授权页面，直接跳转，只能获取用户openid），
							snsapi_userinfo（弹出授权页面，可通过openid拿到昵称、性别、所在地。并且，即使在未关注的情况下，只要用户授权，也能获取其信息）
	state				否	重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
	#wechat_redirect	是	无论直接打开还是做页面302重定向时候，必须带此参数
 */
func GetRedirectURL(w rest.ResponseWriter, r *rest.Request) {
	redirectUrl2 := url.QueryEscape(Redirect_uri)	//对目标url编码
	url2 := fmt.Sprintf(redirectOauthURL,tappid, redirectUrl2,"snsapi_userinfo","237")	//返回实际的url地址
	w.WriteJson(url2)
}


//解码目标地址
func GetCodeforJStoken(w rest.ResponseWriter, r *rest.Request)  {
	w.WriteJson(r)
}