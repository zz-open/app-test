#!/usr/bin/env bash

# 以下代码也可以写成JenkinsFile

HOST_CODE_PATH=./
HOST_CODE_CONF_PATH=./config.yaml

IMAGE_NAME=zz-hello
IMAGE_TAG=v3.0.0

builTagPushHarbor(){
    local SERVER=zz.harbor.com/zzimage
    local SOURCE_IMAGE=${IMAGE_NAME}:${IMAGE_TAG}
    local TARGET_IMAGE=${SERVER}/${IMAGE_NAME}:${IMAGE_TAG}

    echo '123456' | docker login -u admin --password-stdin ${SERVER} && \
    sed -i "s/version:.*/version: \"${IMAGE_TAG}\"/g" ${HOST_CODE_CONF_PATH} && \
    docker build -t ${SOURCE_IMAGE} . && \
    docker tag ${SOURCE_IMAGE} ${TARGET_IMAGE} && \
    docker push ${TARGET_IMAGE}
}

# builTagPushHarbor

builTagPushAli(){
    local SERVER=registry.cn-hangzhou.aliyuncs.com/zzimage
    local SOURCE_IMAGE=${IMAGE_NAME}:${IMAGE_TAG}
    local TARGET_IMAGE=${SERVER}/${IMAGE_NAME}:${IMAGE_TAG}

    # 使用前记得先登录
    # echo 'xxx' | docker login -u xxx --password-stdin registry.cn-hangzhou.aliyuncs.com && \
    sed -i "s/version:.*/version: \"${IMAGE_TAG}\"/g" ${HOST_CODE_CONF_PATH} && \
    docker build -t ${SOURCE_IMAGE} . && \
    docker tag ${SOURCE_IMAGE} ${TARGET_IMAGE} && \
    docker push ${TARGET_IMAGE}
}

builTagPushAli