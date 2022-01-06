# IoT模拟器

## 背景

由于硬件资源的不足,测试和开发时会遇到问题,通过模拟器可以弥补硬件资源不足时的问题

## 安装

`git clone ssh://git@gitlab.yingzi.com:2222/xieke/edge_device_simulator.git`

## 使用方法

本项目使用到了额外的数据库组件, 所有开发和测试都是在PostgreSQL-13.2下完成

项目中使用到的配置变量都是从环境变量中读取

### 部署

#### 通过docker-compoer部署

1. 安装docker. [docker安装教程](https://docs.docker.com/engine/install/)
2. 项目根目录下执行
   `docker-compose up -d`

#### 自行部署

项目中有提供Dockerfile,可以自行build镜像

例如:`docker build -t registry.cn-shenzhen.aliyuncs.com/yingzi/yingzi-edge-device-simulator:v0.0.1 .`

simulator从环境变量中读取,以下为使用到的所有环境变量配置:

```shell
DB_HOST: "simulator_postgresql"
DB_USER: "root"
DB_PASSWORD: "password"
DB_NAME: "postgres"
DB_PORT: "5432"
LOG_LEVEL: "Debug"
LOG_FILENAME: "iot-simulator"
LOG_MAXDAYS: 7
SERVER_PORT: 8088
MAX_OPEN_CONNECTION: 20
MAX_IDLE_CONNECTION: 20
```

### HTTP接口

## 设计文档



