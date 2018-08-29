/**
*使用中控服务器统一获取和刷新Access_token，其他业务逻辑服务器所使用的access_token
均来自于该中控服务器，不应该各自去刷新，否则容易造成冲突，导致access_token覆盖而影响业务；
１、将Access_token存取到redis中
２、定时发送ｇｅｔ请求，更新redis
３、这个get请求是可能失败的 是可能失败的 是可能失败的，虽然这个概率极低，但对于如此重要的参数，
一旦在2个小时的空档期内都无法调用，会产生极其灾难性的后果，别问我怎么知道的

*/
package tokenserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"wechat-tokenServer/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/devfeel/dotweb/framework/convert"
)

type ResAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ErrorCode struct {
	ErrorCode int64  `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

const (
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
)

/**
 *	强制从微信服务器获取accessToken
 *
 *	//调用方法：
 *	acessToken := tokenserver.ResAccessToken{}
 *	acessToken.InitResAccessToken(AppID,AppSecret)
 *
 *	//测试输出
 *	expiresin := strconv.FormatInt(acessToken.ExpiresIn,10)
 *	log.Println(acessToken.AccessToken + " : " + expiresin )
 **/
func (restoken *ResAccessToken) InitResAccessToken(appid, appscret string) (nresAccessToken *ResAccessToken) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, appid, appscret)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("获取accessToken失败：", err)
		logger.Error.Println("获取accessToken失败：", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal("获取accessToken服务器返回异常", AccessTokenURL, response.StatusCode)
		logger.Error.Println("获取accessToken服务器返回异常", AccessTokenURL, response.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		logger.Error.Println(err)
		return
	}
	err = json.Unmarshal(body, restoken) //将json反序列化成struct对象
	if err != nil {
		log.Fatal(err)
		logger.Error.Println(err)
		return
	}
	if restoken.AccessToken == "" {
		errocode := ErrorCode{}
		errocode.SetErrorCode(body)

		//由于微信服务器繁忙导致获取失败，
		//系统休眠5秒后重新获取
		if errocode.ErrorCode == -1 {
			//系统休眠5秒钟
			time.Sleep(5 * time.Second)
			restoken.InitResAccessToken(appid, appscret)
		}
		return
	}

	return restoken
}

/**
 *	返回错误信息赋值到errorCode结构体中
 *	错误时微信会返回错误码等信息，JSON数据包示例如下（该示例为AppID无效错误）:
 *
 * 	{"errcode":40013,"errmsg":"invalid appid"}
 * 	返回码说明
 *
 * 	返回码	说明
 * 	-1	系统繁忙，此时请开发者稍候再试
 *	 0	请求成功
 * 	40001	AppSecret错误或者AppSecret不属于这个公众号，请开发者确认AppSecret的正确性
 * 	40002	请确保grant_type字段值为client_credential
 * 	40164	调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置。（小程序及小游戏调用不要求IP地址在白名单内。）
 */
func (errcode *ErrorCode) SetErrorCode(body []byte) {
	err := json.Unmarshal(body, errcode)
	if err != nil {
		log.Println(err)
		logger.Error.Println(err)
	}
	var strcode string
	switch ecode := errcode.ErrorCode; ecode {
	case -1:
		strcode = "系统繁忙，此时请开发者稍候再试"
	case 0:
		strcode = "请求成功"
	case 40001:
		strcode = "AppSecret错误或者AppSecret不属于这个公众号，请开发者确认AppSecret的正确性"
	case 40002:
		strcode = "请确保grant_type字段值为client_credential"
	case 40164:
		strcode = "调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置。（小程序及小游戏调用不要求IP地址在白名单内。）"
	}

	log.Println("错误代码: " + strconv.FormatInt(errcode.ErrorCode, 10) + ": " + strcode)
	log.Println("错误信息: " + errcode.ErrMsg)
	logger.Error.Println("错误代码: " + strconv.FormatInt(errcode.ErrorCode, 10) + ": " + strcode)
	logger.Error.Println("错误信息: " + errcode.ErrMsg)
}

/**
 *
 *	将accessToken存储到Redis中
 * 	关于过期时间的设置，结合官方文档的说法
 *  ’目前Access_token的有效期通过返回的expire_in来传达，目前是7200秒之内的值。‘
 *	’中控服务器需要根据这个有效时间提前去刷新新access_token。在刷新过程中，中控服务器可对外继续输出的老access_token，‘
 *	’此时公众平台后台会保证在5分钟内，新老access_token都可用，这保证了第三方业务的平滑过渡；‘
 *  中控服务器提前五分钟进行刷新
 */
func (restoken *ResAccessToken) InitRedisDB(redisURL, appid, appSecret string) {

	if restoken.AccessToken != "" {
		//1.链接redis
		c, err := redis.Dial("tcp", redisURL)
		if err != nil {
			log.Println("Redis Connect to redis error", err)
			logger.Error.Println("Redis Connect to redis error", err)
			return
		}
		defer c.Close()

		InvalidTime := restoken.ExpiresIn - 300 //key失效时间

		//2.设置带过期时间的key
		_, errs := c.Do("SET", appid+":"+appSecret, restoken.AccessToken, "EX", InvalidTime)

		if errs != nil {
			log.Println("redis set failed: ", errs)
			logger.Error.Println("redis set failed: ", errs)
		}


		//3.设置定时器，在过期前5分钟从服务器拽取AccessToken，写入到Redis中
		log.Println("正在设置定时器")
		logger.Trace.Println("正在设置定时器")
		ticker := time.NewTicker(time.Duration(InvalidTime) * time.Second)
		log.Println("定时器设置完成，每" + convert.Int642String(InvalidTime) + "秒钟启动从服务器中更新一次token存入Redis中")
		logger.Trace.Println("定时器设置完成，每" + convert.Int642String(InvalidTime) + "秒钟启动从服务器中更新一次token存入Redis中")
		for {
			time := <-ticker.C
			log.Println("正在从服务器获取AccessToken -------")
			logger.Trace.Println("正在从服务器获取AccessToken -------")
			restoken.InitResAccessToken(appid, appSecret) //从服务器拽取AccessToken
			log.Println("已经成功获取AccessToken -------> " + restoken.AccessToken)
			logger.Trace.Println("已经成功获取AccessToken -------> " + restoken.AccessToken)
			log.Println("正在将AccessToken存入Redis中 -------")
			logger.Trace.Println("正在将AccessToken存入Redis中 -------")
			_, errs := c.Do("SET", appid+":"+appSecret, restoken.AccessToken, "EX", InvalidTime) //写入到Redis中
			log.Println("已经成功存入Redis中，key="+appid+":"+appSecret, restoken.AccessToken)
			logger.Trace.Println("已经成功存入Redis中，key="+appid+":"+appSecret, restoken.AccessToken)
			if errs != nil {
				log.Println("redis set failed: ", errs)
				logger.Error.Println("redis set failed: ", errs)
			}
			log.Println("定时器=======>", time.String())
			logger.Trace.Println("定时器=======>", time.String())
		}

	} else {
		log.Println("restoken is null isn't save to Redis DB")
		logger.Error.Println("restoken is null isn't save to Redis DB")
	}

}

/**
从Redis中获取AccessToken值
 */
func (restoken *ResAccessToken) GetAccessTokenFromRedis(redisURL, appid, appSecret string)(naccesstoken *ResAccessToken)  {

	//链接redis
	c, err := redis.Dial("tcp", redisURL)
	if err != nil {
		log.Println("Redis Connect to redis error", err)
		logger.Error.Println("Redis Connect to redis error", err)
		return
	}
	defer c.Close()
	keys := appid+ ":"+appSecret
	is_key_exit, err := redis.Bool(c.Do("EXISTS",keys ))
	if err != nil {
		logger.Error.Println("error:", err)
	}
	logger.Trace.Printf("exists or not: %v \n", is_key_exit)
	accesstokenstr,err := redis.String(c.Do("GET",keys))
	if err != nil{
		logger.Error.Println(err)

	}
	restoken.AccessToken = accesstokenstr

	return restoken
}

//修改redis中的accesstoken值
func (restoken *ResAccessToken) ModifyAccessTokenInRedis(redisURL,appid,appSecret,accesstoken,invalidTime  string)  {
		//1.链接redis
		c, err := redis.Dial("tcp", redisURL)
		if err != nil {
			log.Println("Redis Connect to redis error", err)
			logger.Error.Println("Redis Connect to redis error", err)
			return
		}
		defer c.Close()


		//2.设置带过期时间的key
		_, errs := c.Do("SET", appid+":"+appSecret, accesstoken, "EX", invalidTime)

		if errs != nil {
			log.Println("redis set failed: ", errs)
			logger.Error.Println("redis set failed: ", errs)
		}





}











//**********************************test Support Method **********************************************************/
func Add(a, b int) int {
	c := a + b
	log.Println("Add Invoker" + convert.Int2String(c))
	return c
}

func Routes() {
	http.HandleFunc("/sendjson", SendJSON)
}

func SendJSON(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		Name string
	}{
		Name: "张三",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(u)
}

//测试定时器
func Timer() {
	ticker := time.NewTicker(time.Duration(6900) * time.Second)
	for {
		time := <-ticker.C
		Add(3, time.Second())

		log.Println("定时器=======>", time.String())

	}
}
