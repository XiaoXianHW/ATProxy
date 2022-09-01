package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/XiaoXianHW/ATProxy/config"
	"github.com/XiaoXianHW/ATProxy/console"
	"github.com/XiaoXianHW/ATProxy/service"
	"github.com/XiaoXianHW/ATProxy/version"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	log.SetOutput(color.Output)
	console.SetTitle(fmt.Sprintf("ATProxy %v", version.Version))
	color.HiGreen("感谢使用 ATProxy !\n", )
	color.HiBlue("当前版本 %s !\n", version.Version)
	color.HiBlack("构建Info: %s, %s/%s\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
	go version.CheckUpdate()

	config.LoadConfig()

	for _, s := range config.Config.Services {
		go service.StartNewService(s)
	}
	// hot reload
	// use inotify in linux
	// use Win32 ReadDirectoryChangesW in Windows
	{
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Panic(err)
		}
		defer watcher.Close()
		err = config.MonitorConfig(watcher)
		if err != nil {
			log.Panic("Config Reload Error : ", err)
		}
	}

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
		// stop the program
		// sometimes after the program exits on Windows, the ports are still occupied and "listening".
		// so manually closes these listeners when the program exits.
		for _, listener := range service.ListenerArray {
			if listener != nil { // avoid null pointers
				listener.Close()
			}
		}
	}
}
