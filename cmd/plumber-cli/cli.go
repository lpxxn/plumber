package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/lpxxn/plumber/config"
	"github.com/lpxxn/plumber/src/client"
	"github.com/lpxxn/plumber/src/log"
)

func main() {
	flags := NewFlags()
	flags.Parse(os.Args[1:])
	configFile := flags.Lookup("config").Value.String()
	if configFile == "" {
		log.Errorf("config file is empty")
		return
	}
	_, err := os.Stat(configFile)
	if err != nil {
		log.Error(err)
		return
	}
	cliConf := config.NewCliConf()
	if err := config.ReadFile(configFile, cliConf); err != nil {
		panic(err)
	}
	if err := cliConf.Validate(); err != nil {
		panic(err)
	}
	cli := client.NewClient(cliConf)
	if err := cli.Run(); err != nil {
		panic(err)
	}
	// exit signal
	log.Info("cli ruing....")
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	select {
	case <-cli.GetExitChan():
	case <-ch:
	}
	log.Infof("cli exit!")
}

func NewFlags() *flag.FlagSet {
	flagSet := flag.NewFlagSet("plumber", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")
	return flagSet
}
