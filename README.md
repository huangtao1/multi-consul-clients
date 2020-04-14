# multi-consul-clients  
python、go、java语言客户端使用consul作为注册中心互相访问demo     

## gateway目录  
`spring-cloud-gateway`demo,可以通过gateway来访问其他所有的服务，可以作为前端对接的入口  
修改application.yaml文件可以配置consul相关信息与注册名  

## python目录 
使用`fastapi`构建的python项目 demo  
修改项目中的consul相关配置即可注册到自己的consul上  
`pip install fastapi[all]`  
`uvicorn main:app --host 0.0.0.0`  

## go目录  
使用`gin`构建的go项目 demo  
需要安装gin模块  
`go run main.go`即可启动  

## java目录  
使用`spring-cloud`构建的java项目 demo  
`mvn clean install -DskipTest`  
`java -jar target/**.jar`  
修改application.yaml文件可以配置consul相关信息与注册名  
