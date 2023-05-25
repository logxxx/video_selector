package base

import (
	"flag"
	"github.com/logxxx/utils/log"
	"net"
	"os"
)

var (
	WorkPath = flag.String("work_path", "", "")

	Port          = flag.Int("port", 0, "")
	IsPaused      = false
	ClientVersion = "v1.0.0"
)

func Init() {
	wd, _ := os.Getwd()
	log.Infof("wd:%v", wd)

	flag.Parse()

	if *WorkPath == "" {
		*WorkPath, _ = os.Getwd()
	}
	log.Infof("use work_path:%v", *WorkPath)

	if *Port == 0 {
		uiPort, err := applyFreePort()
		if err != nil {
			panic(err)
		}
		Port = &uiPort
	}
	log.Infof("use Port:%v", *Port)

}

func applyFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	if err != nil {
		return 0, err
	}

	return port, nil
}
