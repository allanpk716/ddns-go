package dns

import (
	"ddns-go/config"
	"ddns-go/util"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	oneJob_s = &util.OneJob{}
	oneJob_c = &util.OneJob{}
	// 缓存的 IP
	storeIPv4 = ""
	// 新查询到的 IP
	queryIPv4 = ""
)

// DNS interface
type DNS interface {
	Init(conf *config.Config) Domains
	// 添加或更新IPV4记录
	AddUpdateIpv4DomainRecords()
	// 添加或更新IPV6记录
	AddUpdateIpv6DomainRecords()
}

// Domains Ipv4/Ipv6 domains
type Domains struct {
	Ipv4Addr    string
	Ipv4Domains []*Domain
	Ipv6Addr    string
	Ipv6Domains []*Domain
}

// Domain 域名实体
type Domain struct {
	DomainName string
	SubDomain  string
	Exist      bool
}

func (d Domain) String() string {
	if d.SubDomain != "" {
		return d.SubDomain + "." + d.DomainName
	}
	return d.DomainName
}

// GetFullDomain 获得全部的，子域名
func (d Domain) GetFullDomain() string {
	if d.SubDomain != "" {
		return d.SubDomain + "." + d.DomainName
	}
	return "@" + "." + d.DomainName
}

// GetSubDomain 获得子域名，为空返回@
// 阿里云，dnspod需要
func (d Domain) GetSubDomain() string {
	if d.SubDomain != "" {
		return d.SubDomain
	}
	return "@"
}

// RunTimer 定时运行
func RunTimer() {
	defer func() {
		if oneJob_s.Running {
			util.CloseFrp(oneJob_s)
		}
		if oneJob_c.Running {
			util.CloseFrp(oneJob_c)
		}

		if err := recover(); err != nil {
			fmt.Println(err)
		}

		log.Printf("Close Frps Frpc Done.")
	}()

	for {
		RunOnce()
		time.Sleep(time.Minute * time.Duration(10))
	}
}

// RunOnce RunOnce
func RunOnce() {
	conf := config.GlobalConfig

	var dnsSelected DNS
	switch conf.DNS.Name {
	case "alidns":
		dnsSelected = &Alidns{}
	case "dnspod":
		dnsSelected = &Dnspod{}
	case "cloudflare":
		dnsSelected = &Cloudflare{}
	case "namesilo":
		dnsSelected = &NamesiloDNS{}
	default:
		dnsSelected = &Alidns{}
	}
	tmpDomains := dnsSelected.Init(&conf)
	dnsSelected.AddUpdateIpv4DomainRecords()
	dnsSelected.AddUpdateIpv6DomainRecords()
	// 获取到的外网 IP
	log.Println("获取到外网的 IP:", tmpDomains.Ipv4Addr)
	// 需要开启 frps 以及 frpc
	nowdir, _ := os.Getwd()
	if initOk := util.InitFrpArgs(nowdir, oneJob_s, oneJob_c); initOk == false {
		log.Println("InitFrpArgs Error.")
		return
	}

	// 第一次
	if storeIPv4 == "" {
		log.Printf("第一次启动 ...")
		storeIPv4 = tmpDomains.Ipv4Addr
		if storeIPv4 == "" {
			log.Println("没有获取到外网 IP")
			return
		}
		log.Println("外网IP:", storeIPv4)
		util.StartFrpThings(oneJob_s, oneJob_c)
	} else {
		// 非第一次
		if tmpDomains.Ipv4Addr == "" {
			log.Println("Try to query Ipv4Addr Error.")
			return
		}
		log.Println("原外网IP:", storeIPv4)
		log.Println("现外网IP:", queryIPv4)
		queryIPv4 = tmpDomains.Ipv4Addr
		if storeIPv4 != queryIPv4 {
			log.Printf("Try ReStart frps frpc ...")
			log.Printf("Close frps ...")
			util.CloseFrp(oneJob_s)
			log.Printf("Close frps Done.")

			log.Printf("Close frpc ...")
			util.CloseFrp(oneJob_c)
			log.Printf("Close frpc Done.")

			// 重新更新缓存的 IP 地址
			storeIPv4 = queryIPv4
			// 先要结束之前运行的 frps 以及 frpc
			util.StartFrpThings(oneJob_s, oneJob_c)
		}
	}
}

// ParseDomain 解析域名
func (domains *Domains) ParseDomain(conf *config.Config) {
	// IPV4
	ipv4Addr := conf.GetIpv4Addr()
	if ipv4Addr != "" {
		domains.Ipv4Addr = ipv4Addr
		domains.Ipv4Domains = parseDomainInner(conf.Ipv4.Domains)
	}

	// IPV6
	ipv6Addr := conf.GetIpv6Addr()
	if ipv6Addr != "" {
		domains.Ipv6Addr = ipv6Addr
		domains.Ipv6Domains = parseDomainInner(conf.Ipv6.Domains)
	}
}

// parseDomainInner 解析域名inner
func parseDomainInner(domainArr []string) (domains []*Domain) {
	for _, domainStr := range domainArr {
		domainStr = strings.Trim(domainStr, " ")
		if domainStr != "" {
			domain := &Domain{}
			sp := strings.Split(domainStr, ".")
			length := len(sp)
			if length <= 1 {
				log.Println(domainStr, "域名不正确")
				continue
			} else if length == 2 {
				domain.DomainName = domainStr
			} else {
				// >=3
				domain.DomainName = sp[length-2] + "." + sp[length-1]
				domain.SubDomain = domainStr[:len(domainStr)-len(domain.DomainName)-1]
			}
			domains = append(domains, domain)
		}
	}
	return
}
