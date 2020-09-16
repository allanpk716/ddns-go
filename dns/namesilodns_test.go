package dns

import (
	"ddns-go/config"
	"fmt"
	"testing"
)

func TestNamesiloDNS(t *testing.T) {

	conf, err := config.GetConfigCache()
	if err != nil {
		return
	}

	dnsSelected := &NamesiloDNS{}
	tmpDomains := dnsSelected.Init(&conf)
	dnsSelected.AddUpdateIpv4DomainRecords()

	fmt.Println("Now IP:", tmpDomains.Ipv4Addr)
}
