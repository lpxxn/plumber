package main

import (
	"flag"
	"os"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/common"
	"github.com/lpxxn/plumber/src/log"
	"github.com/lpxxn/plumber/src/service"
)

func main() {
	flags := NewFlags()
	flags.Parse(os.Args[1:])

	configFile := flags.Lookup("config").Value.String()
	srvConf := config.NewSrvConf()
	if configFile == "" {
		log.Errorf("config file is empty")
		return
	}
	_, err := os.Stat(configFile)
	if err != nil {
		log.Error(err)
		return
	}
	if err := config.ReadFile(configFile, srvConf); err != nil {
		panic(err)
	}
	if err := srvConf.Validate(); err != nil {
		panic(err)
	}
	log.Debugf("config: %+v", *srvConf)
	srv := service.NewService(srvConf)
	wg := &common.WaitGroup{}
	wg.WaitFunc(func() {
		srv.Run()
	})
	wg.Wait()
}

func NewFlags() *flag.FlagSet {
	flagSet := flag.NewFlagSet("plumber", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")
	return flagSet
}
