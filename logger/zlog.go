package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

//日志记录器
var (
	Trace   *log.Logger //记录所有日志
	Info    *log.Logger //记录重要的信息
	Warning *log.Logger //记录需要注意的信息
	Error   *log.Logger //记录非常严重的错误

)

//初始化日志文件 error和trace文件每天创建一个
func init() {

	currpath := GetCurrentPath()
	log.Println("日志路径： " + currpath)
	errfile, err := os.OpenFile(currpath+"/accessTokenServer-errors-"+time.Now().Format("2007-01-02")+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	tracefile, err2 := os.OpenFile(currpath+"/accessTokenServer-Trace-"+time.Now().Format("2007-01-02")+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err2 != nil {
		log.Fatalln("Failed to open Trace log file:", err)
	}

	Trace = log.New(io.MultiWriter(tracefile, os.Stderr), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(errfile, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}


func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	s = strings.Replace(s, "\\", "/", -1)
	s = strings.Replace(s, "\\\\", "/", -1)
	i := strings.LastIndex(s, "/")
	path := string(s[0 : i+1])
	return path
}
