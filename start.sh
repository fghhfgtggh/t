#!/bin/bash

# 后台启动你的服务
./or &

# 等待服务启动
sleep 3

# 使用 nc (netcat) 做端口转发（Render 自带）
while true; do
    nc -l -p ${PORT} -c "nc localhost 6324"
done
