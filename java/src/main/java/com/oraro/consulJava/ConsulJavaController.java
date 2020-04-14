package com.oraro.consulJava;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ConsulJavaController {
    @Autowired
    private ConsulJavaService consulJavaService;

    @GetMapping("/hello_py")
    public String getPythonService() {
        return consulJavaService.helloPy();
    }

    @Autowired
    private ConsulGoService consulGoService;

    @GetMapping("/hello_go")
    public String getGoService() {
        return consulGoService.helloGo();
    }

    @RequestMapping("/java_user")
    public String javaUser() {
        return "this is java";
    }

    @RequestMapping("/health")
    public String health() {
        return "health ok";
    }

}
