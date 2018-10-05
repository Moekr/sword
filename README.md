Sword
=====
仿SmokePing的多点ping值监测工具 https://ping.moekr.com

### 编译安装

1. 确保有Go环境
2. make，交叉编译则加上目标系统和架构的环境变量如`GOARCH=amd64 GOOS=linux make`
3. 将output文件夹拷贝至目标机器
4. (按需)修改script和service中的路径

### 启动参数

```
服务端：
sword -s -b 127.0.0.1:7901 -c /opt/sword/conf.json -d /opt/sword/data/ -t token

客户端：
sword -u http://localhost:7901 -i 1 -t token
```

其中：

- -s：申明以服务端模式启动，默认为false即以客户端模式启动
- -b：监听的地址和端口，默认`0.0.0.0:7901`，服务端专用
- -c：配置文件位置，默认`./conf.json`，服务端专用，配置文件模板见conf.sample.json
- -d：数据文件目录，默认`./data/`，服务端专用
- -u：服务端通信地址，默认`http://localhost:7901`，客户端专用
- -i：客户端编号，默认-1，客户端专用，该编号必须对应服务端配置文件中的一个observer id
- -t：HTTP Token，用于客户端与服务端通信时的鉴权，默认为空即不鉴权，服务端与客户端该参数必须一致

此外：

- -v：以debug模式启动，会输出较多日志

### 其他

1. Debian/Ubuntu可以将sword.server.service/sword.client.service注册为系统服务
2. 修改配置文件后推荐使用script/reconf.sh重启(Debian/Ubuntu)
3. 每一小时将生成一个备份文件，可以设定定时任务调用script/rmbak.sh清理备份文件

### 协议

GPLv3
