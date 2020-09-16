package util

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	oneJob_s = &OneJob{}
	oneJob_c = &OneJob{}
	// 缓存的 IP
	storeIPv4 = ""
	// 新查询到的 IP
	queryIPv4 = ""
	nowdir    = ""
	ip_1      = "1.1.1.1"
	ip_2      = "2.2.2.2"

	index = 0
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, string(os.PathSeparator)))
}

func TestStartFrpThings(t *testing.T) {
	defer func() {
		if oneJob_s.Running {
			CloseFrp(oneJob_s)
		}
		if oneJob_c.Running {
			CloseFrp(oneJob_c)
		}

		if err := recover(); err != nil {
			// 第三个是人为制造的故障
			if index != 3 {
				t.Fatal(err)
			}
		}

		fmt.Println("Close Frps Frpc Done.")
	}()

	// 需要开启 frps 以及 frpc

	orgIP := ""
	newIP := ""

	for {
		if index == 4 {
			break
		} else if index == 0 {
			// 第一次，正常启动
			orgIP = ip_1
			newIP = ip_2
			nowdir, _ = os.Getwd()
			nowdir = getParentDirectory(nowdir)
		} else if index == 1 {
			// 第二次，期望跳过 restart frp
			orgIP = ip_1
			newIP = ip_1
			nowdir, _ = os.Getwd()
			nowdir = getParentDirectory(nowdir)
		} else if index == 2 {
			// 第三次，期望 restart frp
			orgIP = ip_1
			newIP = ip_2
			nowdir, _ = os.Getwd()
			nowdir = getParentDirectory(nowdir)
		} else if index == 3 {
			// 第四次，期望 InitFrpArgs == false
			// 这里没有能够进入项目的 root 目录，所以人为制造 frps 文件不存在的故障
			orgIP = ip_1
			newIP = ip_2
			nowdir, _ = os.Getwd()
		}

		// 第四次，期望 InitFrpArgs == false
		if index == 3 {
			if initOk := InitFrpArgs(nowdir, oneJob_s, oneJob_c); initOk == true {
				panic("InitFrpArgs Error.")
			}
			index++
			continue
		} else {
			if initOk := InitFrpArgs(nowdir, oneJob_s, oneJob_c); initOk == false {
				panic("InitFrpArgs Error.")
			}
		}

		// 第一次
		if storeIPv4 == "" {
			storeIPv4 = orgIP
			if StartFrpThings(oneJob_s, oneJob_c) == false {
				panic("Start frps frpc Error")
			}
		} else {
			// 第二、第三次
			queryIPv4 = newIP
			if storeIPv4 != queryIPv4 {
				fmt.Println("Try ReStart frps frpc ...")

				if oneJob_s.Running {
					fmt.Println("Close frps ...")
					if CloseFrp(oneJob_s) == false {
						panic("Close frps Error")
					}
				}
				fmt.Println("frps is Closed.")

				if oneJob_c.Running {
					fmt.Println("Close frpc ...")
					if CloseFrp(oneJob_c) == false {
						panic("Close frpc Error")
					}
				}
				fmt.Println("frpc is Closed.")

				// 先要结束之前运行的 frps 以及 frpc
				if StartFrpThings(oneJob_s, oneJob_c) == false {
					panic("Start frps frpc Error")
				}
			}
		}

		fmt.Println("Index -- ", index, "Done.")
		index++
	}
}
