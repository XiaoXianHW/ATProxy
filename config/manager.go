package config

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/XiaoXianHW/ATProxy/common/set"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (
	Config     configMain
	Lists      map[string]*set.StringSet
	reloadLock sync.Mutex
)

func LoadConfig() {
	configFile, err := os.ReadFile("ATProxy.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("配置文件不存在, 已自动生成新的配置文件!")
			generateDefaultConfig()
			goto success
		} else {
			log.Panic(color.HiRedString("加载配置时出现意外错误: %s", err.Error()))
		}
	}

	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		log.Panic(color.HiRedString("配置格式错误: %s", err.Error()))
	}

success:
	LoadLists(false)
	log.Println(color.HiYellowString("已成功载入配置文件!"))
}

func generateDefaultConfig() {
	file, err := os.Create("ATProxy.json")
	if err != nil {
		log.Panic("创建配置文件失败：", err.Error())
	}
	Config = configMain{
		Services: []*ConfigProxyService{
			{
				Name:          "Hypixel",
				TargetAddress: "mc.hypixel.net",
				TargetPort:    25565,
				Listen:        25565,
				Flow:          "auto",
				Minecraft: minecraft{
					EnableHostnameRewrite: true,
					IgnoreFMLSuffix:       true,
					OnlineCount: onlineCount{
						Max:            10,
						Online:         -1,
						EnableMaxLimit: false,
					},
					MotdFavicon:     "{DEFAULT_MOTD}",
					MotdDescription: "§e ATProxy 代理已正常启动",
				},
			},
		},
		Lists: map[string][]string{
			//"test": {"foo", "bar"},
		},
	}
	newConfig, _ :=
		json.MarshalIndent(Config, "", "    ")
	_, err = file.WriteString(strings.ReplaceAll(string(newConfig), "\n", "\r\n"))
	file.Close()
	if err != nil {
		log.Panic("保存配置文件失败:", err.Error())
	}
}

func LoadLists(isReload bool) bool {
	reloadLock.Lock()
	if isReload {
		configFile, err := os.ReadFile("ATProxy.json")
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(color.HiRedString("重新加载失败：配置文件不存在。"))
			} else {
				log.Println(color.HiRedString("重新加载配置时出现意外错误: %s", err.Error()))
			}
			reloadLock.Unlock()
			return false
		}

		err = json.Unmarshal(configFile, &Config)
		if err != nil {
			log.Println(color.HiRedString("无法重新加载：配置格式错误: %s", err.Error()))
			reloadLock.Unlock()
			return false
		}
	}
	//log.Println("Lists:", Config.Lists)
	if l := len(Config.Lists); l == 0 { // if nothing in Lists
		Lists = map[string]*set.StringSet{} // empty map
	} else {
		Lists = make(map[string]*set.StringSet, l) // map size init
		for k, v := range Config.Lists {
			//log.Println("List: Loading", k, "value:", v)
			set := set.NewStringSetFromSlice(v)
			Lists[k] = &set
		}
	}
	Config.Lists = nil // free memory
	reloadLock.Unlock()
	runtime.GC()
	return true
}

func MonitorConfig(watcher *fsnotify.Watcher) error {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write { // config reload
					log.Println(color.HiMagentaString("热重载：检测到ID名单文件更改,正在重新加载..."))
					if LoadLists(true) { // reload success
						log.Println(color.HiMagentaString("热重载：成功重新加载ID名单."))
					} else {
						log.Println(color.HiMagentaString("热重载：无法重新加载ID名单."))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				log.Println(color.HiRedString("热重载错误 : ", err))
			}
		}
	}()

	return watcher.Add("ATProxy.json")
}
