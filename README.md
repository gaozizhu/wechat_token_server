#### 本项目是采用golang语言构建的，用于微信公众号的AccessToken中控服务器

##### 设计思路：

1、对外提供基于RestAPI和Redis数据库的AccessToken服务。
2、将服务器和Redis封装到Docker中，方便部署。
3、考虑到服务器的健壮性和稳定性，采用redis集群多物理服务器部署并且和Mysql数据库进行整合。

##### 初始化：

1、从微信服务器获取AccessToken。

2、将AccessToken存储到Redis中。



##### 当前进展

+ 由于从微信服务器中Get请求AccessToken，这个GET是有可能失败的，虽然这个概率极低，但对于如此重要的参数，一旦在2个小时的空档期内都无法调用，会产生灾难性的后果。

  + 1、考虑将中控器在阿里和腾讯以及其他自建服务器等多个服务器上部署，通过网关进行跳转，一旦其中一台无法访问立即启动另一台服务器上的中控服务。
  + 2、提供手动设置AccessToken的接口，应急处理时手动将其他渠道获得的AccessToken设置到Redis数据库中。
  + 3、手动设置支持命令行方式

+ go的函数递归调用

  + 闭包的含义和用法

    f := function(x,y int) int{

    .......

    }

  + 递归调用的含义和方法

+ 使用redis的key过期事件

  + 上网查了redis的key过期时间通知有时候丢失的现象。研究一下解决办法或替换方案。
    + 消息中间件
    + redis LAU语言
    + go的定时任务

+ redis官方中文文档 http://redisdoc.com/index.html

+ golang 单元测试

  + 测试文件和要测试的.go文件放到同一个目录下。

  + 测试文件名称 被测试文件名_test.go

  + 测试方法名称 

    ``` go
    //测试函数的签名必须接收一个指向testing.T类型的指针，并且不能返回任何值
    func TestAdd(t *testing.T)  {
    	sum := Add(1,2)
    	if sum == 3 {
    		t.Log("the result is ok")
    
    	}else{
    		t.Fatal("the result is wrong")
    	}
    }
    ```

  + 运行测试 进入测试文件目录 

    ``` go
    go test -v
    ```

    

  + 针对模拟网络访问，标准库了提供了一个httptest包，可以让我们模拟http的网络调用，下面举个例子了解使用。

    首先我们创建一个处理HTTP请求的函数，并注册路由

    ``` go
    package common
    
    import (
    	"net/http"
    	"encoding/json"
    )
    
    func Routes(){
    	http.HandleFunc("/sendjson",SendJSON)
    }
    
    func SendJSON(rw http.ResponseWriter,r *http.Request){
    	u := struct {
    		Name string
    	}{
    		Name:"张三",
    	}
    
    	rw.Header().Set("Content-Type","application/json")
    	rw.WriteHeader(http.StatusOK)
    	json.NewEncoder(rw).Encode(u)
    }
    ```

    ```go
    //测试单个文件，一定要带上被测试的原文件
    
        go test -v  wechat_test.go wechat.go 
    
       
    //测试单个方法
    
        go test -v -test.run TestRefreshAccessToken
    ```

  安装编译运行
 ```go
//引入包：
go get github.com/gomodule/redigo/redis
go get -u github.com/devfeel/dotweb
go get github.com/ant0ine/go-json-rest/rest 
//rest相关文档 https://github.com/ant0ine/go-json-rest

 ```

用curl测试post方法

curl -l -H "Content-type: application/json" -X POST -d '{"token":"13521389587"}' http://localhost:8777/accesstoken/adf

github
1-添加远程地址

```shell
git remote add origin https://github.com/gaozizhu/wechat_token_server.git
```

$ git config branch.master.remote origin  
$ git config branch.master.merge refs/heads/master 
 2-提交上传本地文件

git add . 			//将当前全部文件添加到“暂存区”

git commit -m "init"	//提交到暂存区

