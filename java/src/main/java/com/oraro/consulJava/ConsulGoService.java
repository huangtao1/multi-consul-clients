package com.oraro.consulJava;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.RequestMapping;


@FeignClient(value = "consul-go", fallback = FeignFallBack.class)
public interface ConsulGoService {

    @RequestMapping("/go_user")
    String helloGo();
}

