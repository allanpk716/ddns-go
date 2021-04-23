package main

import (
	"ddns-go/config"
	"ddns-go/dns"
	"ddns-go/static"
	"ddns-go/web"
	"log"
	"net/http"
	"time"
)

const port = "9876"

func main() {
	// 启动静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(static.AssetFile())))
	http.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(static.AssetFile())))

	http.HandleFunc("/", config.BasicAuth(web.Writing))
	http.HandleFunc("/save", config.BasicAuth(web.Save))
	http.HandleFunc("/logs", config.BasicAuth(web.Logs))

	// 打开浏览器
	// go util.OpenExplorer("http://127.0.0.1:" + port)
	log.Println("启动端口", port, "...")

	// 提前读取一次
	config.GetConfigCache()

	// 定时运行
	go dns.RunTimer()

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Println("启动端口发生异常, 1分钟后自动关闭此窗口", err)
		time.Sleep(time.Minute)
	}
}
