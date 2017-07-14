package app

import (
	"github.com/zhiyuan1024/sysinfo/app/cpuinfo"
	"log"
	"os"
)

type SysInfo interface {
	Collecting(line chan string)
}

type ApplicationConfig struct {
	DataFile string
}

type Application struct {
	AppConfig ApplicationConfig
	cpu       SysInfo
	lineChan  chan string
}

func (app *Application) Start() {
	go app.cpu.Collecting(app.lineChan)
	for {
		select {
		case line, ok := <-app.lineChan:
			if !ok {
				log.Fatalf("catche sysinfo error, chan close")
			}
			if err := app.write(line); err != nil {
				log.Fatalf("write sysinfo to %s error, err = %v", app.AppConfig.DataFile, err)
			}
		}
	}
}

func NewApplication() *Application {
	app := new(Application)
	app.AppConfig.DataFile = "data/sysinfo.data"
	app.lineChan = make(chan string)
	app.cpu = cpuinfo.NewCPUInfo()
	return app
}

func (app *Application) write(line string) error {
	f, err := os.OpenFile(app.AppConfig.DataFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write([]byte(line)); err != nil {
		return err
	}
	return nil
}
