package com.oraro.consulJava;

import org.springframework.stereotype.Component;

@Component
public class FeignFallBack implements ConsulJavaService {
    public String helloPy() {
        return null;
    }
}
