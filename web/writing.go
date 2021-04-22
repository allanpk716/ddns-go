package web

import (
	"ddns-go/config"
	"ddns-go/util"
	"log"
	"strings"

	"fmt"
	"html/template"
	"net/http"
)

// Writing 填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	tempPath, err := util.GetStaticResourcePath("static/pages/writing.html")
	if err != nil {
		log.Println("Asset was not found.")
		return
	}
	tmpl, err := template.ParseFiles(tempPath)
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	conf, err := config.GetConfigCache()
	if err == nil {
		config.GlobalConfig = conf
		// 已存在配置文件，隐藏真实的ID、Secret
		idHide, secretHide := getHideIDSecret(&conf)
		conf.DNS.ID = idHide
		conf.DNS.Secret = secretHide
		tmpl.Execute(writer, &conf)
		return
	}

	// 默认值
	if conf.Ipv4.URL == "" {
		conf.Ipv4.URL = "https://api-ipv4.ip.sb/ip"
		conf.Ipv4.Enable = true
	}
	if conf.Ipv6.URL == "" {
		conf.Ipv6.URL = "https://api-ipv6.ip.sb/ip"
	}
	if conf.DNS.Name == "" {
		conf.DNS.Name = "alidns"
	}

	tmpl.Execute(writer, conf)
}

// 显示的数量
const displayCount int = 3

// hideIDSecret 隐藏真实的ID、Secret
func getHideIDSecret(conf *config.Config) (idHide string, secretHide string) {
	if len(conf.DNS.ID) > displayCount {
		idHide = conf.DNS.ID[:displayCount] + strings.Repeat("*", len(conf.DNS.ID)-displayCount)
	}
	if len(conf.DNS.Secret) > displayCount {
		secretHide = conf.DNS.Secret[:displayCount] + strings.Repeat("*", len(conf.DNS.Secret)-displayCount)
	}
	return
}
