package version

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func printErr(err error) {
	log.Printf("检查更新时发生错误, 原因: %v.", err.Error())
	log.Println(`你可以尝试在这里查看最新版本 https://github.com/XiaoXianHW/ATProxy`)
}

func CheckUpdate() {
	resp, err := http.Get(`https://htmlbackup-1302685493.cos.ap-hongkong.myqcloud.com/ATProxy/version.go`)
	if err != nil {
		printErr(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printErr(err)
		return
	}
	if strings.Contains(string(body), Version) {
		fmt.Println("你当前的ATProxy为最新版本!")
	} else {
		fmt.Println("你当前的ATProxy版本为旧版本，请前往 https://github.com/XiaoXianHW/ATProxy/releases 获取最新版本！")
	}
}
