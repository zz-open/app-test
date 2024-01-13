#!/usr/bin/env bash

docker rmi registry.cn-hangzhou.aliyuncs.com/zzimage/zz-hello:v1.0.0
docker rmi zz-hello:v1.0.0

docker rmi registry.cn-hangzhou.aliyuncs.com/zzimage/zz-hello:v2.0.0
docker rmi zz-hello:v2.0.0

docker rmi registry.cn-hangzhou.aliyuncs.com/zzimage/zz-hello:v3.0.0
docker rmi zz-hello:v3.0.0

docker image prune -f

docker images