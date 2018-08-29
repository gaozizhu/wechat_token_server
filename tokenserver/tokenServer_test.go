package tokenserver

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wechat-tokenServer/logger"
)

func init() {
	Routes()
}

func TestSendJSON(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/sendjson", nil)
	if err != nil {
		t.Fatal("创建Request失败")
	}
	rw := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rw, req)

	log.Println("code:", rw.Code)
	log.Println("body:", rw.Body.String())
	logger.Trace.Println("body:", rw.Body.String())
}

func TestAdd(t *testing.T) {
	sum := Add(1, 2)
	if sum == 3 {
		t.Log("the result is ok")

	} else {
		t.Fatal("the result is wrong")
	}
}

//测试定时器
func TestTimer(t *testing.T) {
	Timer()
}
