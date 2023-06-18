package main

import (
	"context"
	"github.com/kennnyz/unixGo/configs"
	"github.com/kennnyz/unixGo/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	configPath = "D:\\unixGoWB\\unixSocketsGO\\configs\\config.json"
)

func main() {
	var wg sync.WaitGroup
	cfg := configs.ReadConfig(configPath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	UnixServer := server.NewServer(cfg.ListenAddress)
	err := UnixServer.Start(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	// Cleanup the socket file.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-c
		err := UnixServer.Close()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Remove(cfg.ListenAddress)
		os.Exit(0)
	}()
	wg.Wait()
}
