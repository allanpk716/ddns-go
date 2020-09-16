<a href="https://github.com/jeessy2/ddns-go/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/jeessy2/ddns-go.svg?logo=github&style=flat-square"></a>

# ddns-go
- 自动获得你的公网IPV4或IPV6并解析到域名中
- 支持Mac、Windows、Linux系统，支持ARM、x86架构
- 支持的域名服务商 `Alidns(阿里云)` `Dnspod(腾讯云)` `Cloudflare`  NameSilo
- 间隔5分钟同步一次
- 支持多个域名同时解析，公司必备
- 支持多级域名
- 网页中配置，简单又方便，可设置登录用户名和密码
- 网页中方便快速查看最近50条日志，不需要跑docker中查看

## 系统中使用
- 下载并解压[https://github.com/jeessy2/ddns-go/releases](https://github.com/jeessy2/ddns-go/releases)
- 双击运行，程序自动打开[http://127.0.0.1:9876](http://127.0.0.1:9876)，修改你的配置，成功
- [可选] 加入到开机启动中，需自行搜索

## Docker中使用
```
docker run -d \
  --name ddns-go \
  --restart=always \
  -p 9876:9876 \
  jeessy/ddns-go
```
- 在网页中打开`http://主机IP:9876`，修改你的配置，成功
- [可选] docker中默认不支持ipv6，需自行探索如何开启

## 使用IPV6
- 前提：你的电脑或终端能正常获取IPV6
- Windows/Mac系统推荐在 `系统中使用`，Windows/Mac桌面版的docker不支持`--net=host`
- Linux的x86或arm架构，如服务器、群晖、xx盒子等等，推荐使用`--net=host`模式，简单点
  ```
  docker run -d \
    --name ddns-go \
    --restart=always \
    --net=host \
    jeessy/ddns-go
  ```
- [可选] 使用IPV6后，建议设置登录用户名和密码

![avatar](ddns-web.png)

## 扩展功能，使用 FRP

如果家里有动态的公网 IP，那么很可能想要能使用这个优势搭建自己的代理、内网穿透， FRP 是一个很好的选择。但是就目前的测试，FRP 在动态网外 IP 的情景下，无法能够在 IP 变动后正常工作。本功能仅简单实现了外网 IP 变动后， FRP 的自动重启。

下面是具体的使用设置：

1. 请**自学** FRP 相关配置知识
2. 在 ddns-go 程序根目录下，新建 frpThings 文件夹
3. 放入 frps、frps.ini、frpc、frpc.ini ，务必保证名称（这个也得看你是啥系统，Widnows 就是 frps.exe 其他系统无需 .exe 后缀，以此类推）一样，否则无法启动
4. 正常启动 ddns-go 即可，会跟随公网 IP 检测逻辑执行判断是否需要重启

如果需要在 Docker 中使用，也需要创建对应的映射目录，同时放入  frps、frps.ini、frpc、frpc.ini 

```
localFolder/frpThings:/app/frpThings
```

## Development

```
go get -u github.com/go-bindata/go-bindata/...
go-bindata -debug -pkg util -o util/staticPagesData.go static/pages/...
go-bindata -pkg static -ignore js_css_data.go -o static/js_css_data.go -fs -prefix "static/" static/
```

## Release
```
go-bindata -pkg util -o util/staticPagesData.go static/pages/...
go-bindata -pkg static -ignore js_css_data.go -o static/js_css_data.go -fs -prefix "static/" static/

# 自动发布
git tag v0.0.x -m "xxx" 
git push --tags
```