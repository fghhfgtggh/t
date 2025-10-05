package main

import (
    "io"
    "log"
    "net"
    "os"
    "os/exec"
    "time"
)

func main() {
    // 启动你的 or 二进制文件
    cmd := exec.Command("./or")
    if err := cmd.Start(); err != nil {
        log.Fatal("Failed to start or:", err)
    }
    log.Println("Started ./or")
    
    // 等待服务启动
    time.Sleep(3 * time.Second)
    
    // 获取 Render 分配的端口
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }
    
    // 监听 Render 的端口
    listener, err := net.Listen("tcp", "0.0.0.0:"+port)
    if err != nil {
        log.Fatal("Failed to listen:", err)
    }
    
    log.Printf("Proxy listening on :%s -> forwarding to localhost:6324", port)
    
    // 接受并转发连接
    for {
        client, err := listener.Accept()
        if err != nil {
            log.Println("Accept error:", err)
            continue
        }
        go handleConnection(client)
    }
}

func handleConnection(client net.Conn) {
    defer client.Close()
    
    // 连接到你的实际服务（6324端口）
    backend, err := net.Dial("tcp", "127.0.0.1:6324")
    if err != nil {
        log.Println("Backend connection error:", err)
        return
    }
    defer backend.Close()
    
    // 双向数据转发
    done := make(chan struct{}, 2)
    
    go func() {
        io.Copy(backend, client)
        done <- struct{}{}
    }()
    
    go func() {
        io.Copy(client, backend)
        done <- struct{}{}
    }()
    
    <-done
}
