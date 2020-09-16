package dns

import (
	"ddns-go/config"
	"log"

	"github.com/nrdcg/namesilo"
	namesiloSDK "github.com/nrdcg/namesilo"
)

// Alidns 阿里云dns实现
type NamesiloDNS struct {
	client *namesilo.Client
	Domains
}

// Init 初始化
func (namesiloDNS *NamesiloDNS) Init(conf *config.Config) Domains {
	transport, err := namesiloSDK.NewTokenTransport(conf.DNS.Secret)
	if err != nil {
		log.Fatal(err)
	}

	client := namesilo.NewClient(transport.Client())

	namesiloDNS.client = client

	namesiloDNS.Domains.ParseDomain(conf)

	return namesiloDNS.Domains
}

// AddUpdateIpv4DomainRecords 添加或更新IPV4记录
func (namesiloDNS *NamesiloDNS) AddUpdateIpv4DomainRecords() {
	namesiloDNS.addUpdateDomainRecords("A")
}

// AddUpdateIpv6DomainRecords 添加或更新IPV6记录
func (namesiloDNS *NamesiloDNS) AddUpdateIpv6DomainRecords() {
	namesiloDNS.addUpdateDomainRecords("AAAA")
}

func (namesiloDNS *NamesiloDNS) addUpdateDomainRecords(recordType string) {
	ipAddr := namesiloDNS.Ipv4Addr
	domains := namesiloDNS.Ipv4Domains
	if recordType == "AAAA" {
		ipAddr = namesiloDNS.Ipv6Addr
		domains = namesiloDNS.Ipv6Domains
	}

	if ipAddr == "" {
		return
	}

	for _, domain := range domains {
		// 获取已经存在的列表
		params_ListRecord := &namesiloSDK.DnsListRecordsParams{}
		params_ListRecord.Domain = domain.DomainName
		result_ListRecord, err := namesiloDNS.client.DnsListRecords(params_ListRecord)
		if err != nil {
			log.Printf("namesilo DnsListRecords Error: %s", err)
			return
		}
		if result_ListRecord == nil {
			log.Printf("namesilo DnsListRecords result is nil: %s", domain)
			return
		}
		if result_ListRecord.Reply.Code != "300" {
			log.Printf("namesilo DnsListRecords require code is not 300")
			return
		}
		// 域名是否存在
		bfound := false
		id := ""
		for _, oneRecord := range result_ListRecord.Reply.ResourceRecord {
			if oneRecord.Host == domain.SubDomain+"."+domain.DomainName {
				id = oneRecord.RecordID
				bfound = true
				break
			}
		}
		if bfound == true {
			// 存在则更新
			namesiloDNS.modify(id, domain, recordType, ipAddr)
		} else {
			// 不存在则新建
			namesiloDNS.create(domain, recordType, ipAddr)
		}
	}
}

// 新增
func (namesiloDNS *NamesiloDNS) create(domain *Domain, recordType string, ipAddr string) {
	params_DnsAddRecord := &namesiloSDK.DnsAddRecordParams{}
	params_DnsAddRecord.Domain = domain.DomainName
	params_DnsAddRecord.Type = recordType
	params_DnsAddRecord.Host = domain.SubDomain
	params_DnsAddRecord.Value = ipAddr
	params_DnsAddRecord.Distance = 0
	params_DnsAddRecord.TTL = 3600

	result_DnsAddRecord, err := namesiloDNS.client.DnsAddRecord(params_DnsAddRecord)
	if err != nil {
		log.Printf("namesilo DnsAddRecord Error: %s", err)
		return
	}
	if result_DnsAddRecord == nil {
		log.Printf("namesilo DnsAddRecord result is nil: %s", domain)
		return
	}
}

// 更新
func (namesiloDNS *NamesiloDNS) modify(id string, domain *Domain, recordType string, ipAddr string) {
	params_DnsUpdateRecord := &namesiloSDK.DnsUpdateRecordParams{}
	params_DnsUpdateRecord.Domain = domain.DomainName
	params_DnsUpdateRecord.ID = id
	params_DnsUpdateRecord.Host = domain.SubDomain
	params_DnsUpdateRecord.Value = ipAddr
	params_DnsUpdateRecord.Distance = 0
	params_DnsUpdateRecord.TTL = 3600

	result_DnsUpdateRecord, err := namesiloDNS.client.DnsUpdateRecord(params_DnsUpdateRecord)
	if err != nil {
		log.Printf("namesilo DnsUpdateRecord Error: %s", err)
		return
	}
	if result_DnsUpdateRecord == nil {
		log.Printf("namesilo DnsUpdateRecord result is nil: %s", domain)
		return
	}
}
