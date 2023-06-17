package main

import (
	"github.com/kennnyz/unixGo/configs"
	"github.com/kennnyz/unixGo/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	configPath = "configs/config.json"
)

func main() {
	var wg sync.WaitGroup
	cfg := configs.ReadConfig(configPath)

	server := server.NewServer(cfg.ListenAddress)
	err := server.Start()
	if err != nil {
		log.Println(err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range server.MsgChan {
			log.Printf("%s: %s\n", msg.From, msg.Message)
		}
	}()

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-c
		os.Remove(cfg.ListenAddress)
		os.Exit(1)
	}()
	wg.Wait()
}
