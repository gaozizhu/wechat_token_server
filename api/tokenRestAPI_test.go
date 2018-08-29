package api

import (
	"testing"
	"log"
	"strings"
)

func TestGetAccessTokenFromRedis(t *testing.T) {
	//acessToken := tokenserver.ResAccessToken{}
	//acessToken.GetAccessTokenFromRedis("127.0.0.1:6379", "wx10a300ee87ff239f", "845c55f125af95ae5045b974bc3ed615")
	//jace ,err := json.Marshal(acessToken)
	//if err != nil{
	//	logger.Error.Println(err)
	//}
	//log.Println(string(jace))

	parts := strings.Split("accesstoken/:ad", "/")
	log.Println(parts)
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}

}