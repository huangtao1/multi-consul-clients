package com.oraro.consulJava;

import org.springframework.stereotype.Component;

@Component
public class GoFeignFallBack implements ConsulGoService {
    public String helloGo() {
        return null;
    }
}
