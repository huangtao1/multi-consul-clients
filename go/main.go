package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var serviceId string

func main() {
	router := gin.Default()
	//初始化连接
	client := initConsul()
	//健康检查接口
	router.GET("/check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})
	router.GET("/go_user", func(context *gin.Context) {
		context.JSON(http.StatusOK,gin.H{"message":"this is go"})
	})
	//调用python
	router.GET("/hello_py", func(context *gin.Context) {
		s := connectService("consul-py", "python_user", client)
		context.String(http.StatusOK, s)
	})
	//调用java
	router.GET("/hello_java", func(context *gin.Context) {
		s := connectService("consul-java", "java_user", client)
		context.String(http.StatusOK, s)
	})
	//向consul注册自己
	registerConsul(client)
	//启动监听
	srv := &http.Server{
		Addr:    ":8777",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	//监听退出操作
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	//consul deregister
	err := client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		log.Fatal("Deregister failed:", err)
	}
	log.Println("解除注册成功!!!!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func initConsul() *api.Client {
	config := api.DefaultConfig()
	//定义一个consul client地址
	config.Address = "192.168.9.50:8500"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("connect consul failed:" + err.Error())
	}
	return client
}

//注册consul
func registerConsul(client *api.Client) {
	serviceId = "consul-go" + "-" + uuid.NewV4().String()
	register := api.AgentServiceRegistration{}
	register.Name = "consul-go"
	register.Port = 8777
	register.Address = "192.168.9.52"
	register.ID = serviceId
	//定义心跳检查
	check := api.AgentServiceCheck{}
	check.TCP = fmt.Sprintf("%s:%d", register.Address, register.Port)
	check.Interval = "3s"
	check.DeregisterCriticalServiceAfter = "1m"
	register.Check = &check
	//向consul注册
	err := client.Agent().ServiceRegister(&register)
	if err != nil {
		log.Fatal("register consul failed:", err.Error())
	}
}

//解除注册
func deregisterService(client *api.Client, serviceId string) {
	err := client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		log.Fatal("deregister consul failed:", err.Error())
	}
}
func Get(url string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

//调用其他服务
func connectService(serviceName, requestUrl string, client *api.Client) (result string) {
	queryOption := api.QueryOptions{}
	services, _, err := client.Catalog().Service(serviceName, "", &queryOption)
	if err != nil {
		return "ERROR:" + err.Error()
	}
	//随机获取一个实例
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(services))
	serviceInstance := services[index]
	serviceUrl := fmt.Sprintf("http://%s:%d/%s", serviceInstance.ServiceAddress, serviceInstance.ServicePort, requestUrl)
	log.Println(serviceUrl)
	//发起请求
	result = Get(serviceUrl)
	return result
}
