from fastapi import FastAPI
from fastapi.logger import logger
import consul
import random
import requests
import uuid

app = FastAPI()
service_name = "consul-py"
service_id = service_name + '-' + str(uuid.uuid1())
# 实例consul类
c = consul.Consul(host="192.168.9.50", port=8500)
# 向server注册
service = c.agent.service
service.register(name=service_name, address="192.168.9.52", port=8000, service_id=service_id,
                 check=consul.Check().tcp(host="192.168.9.52", port=8000, interval="3s", deregister="1m"))


def get_other_instance_result(service_name, prefix_url):
    # 调用consul api
    consul_service = c.catalog.service(service_name)
    print(consul_service)
    consul_service_instances = []
    # 获取实例的真实ip
    if len(consul_service[1]) > 0:
        instance_address_infos = consul_service[1]
        for instance_address_info in instance_address_infos:
            consul_service_instances.append(
                "http://" + instance_address_info["ServiceTaggedAddresses"]["lan_ipv4"]["Address"] + ":" +
                str(instance_address_info["ServiceTaggedAddresses"]["lan_ipv4"]["Port"]) + "/")
    else:
        return "no instance"
    # 随机返回一个可用实例
    instance_url = consul_service_instances[random.randint(0, len(consul_service_instances) - 1)]
    print(instance_url)
    try:
        # 发起请求
        result = requests.get(instance_url + prefix_url).text
    except Exception as e:
        result = "ERROR:" + str(e)
    return result


@app.get("/python_user")
async def hello():
    return {"message": "this is python"}


# 调用go
@app.get("/hello_go")
async def hello_go():
    result = get_other_instance_result("consul-go", "go_user")
    if result == "no instance" or result.startswith("ERROR:"):
        return {"message": "调用失败:%s" % result}
    else:
        return {"message": result}


# 调用java
@app.get("/hello_java")
async def hello_java():
    result = get_other_instance_result("consul-java", "java_user")
    if result == "no instance" or result.startswith("ERROR:"):
        return {"message": "调用失败:%s" % result}
    else:
        return {"message": result}


# 关闭前自动解注册
@app.on_event("shutdown")
async def shutdown_event():
    service.deregister(service_id)
    logger.info("解注册成功!!")
