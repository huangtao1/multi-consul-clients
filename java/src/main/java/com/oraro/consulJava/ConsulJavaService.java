package com.oraro.consulJava;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.RequestMapping;


@FeignClient(value = "consul-py", fallback = FeignFallBack.class)
public interface ConsulJavaService {
    @RequestMapping("/python_user")
    String helloPy();
}

