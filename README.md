# unixSocketsGO
unixSocketsGO

Библиотека для взаимодействие через unix сокеты.

Необходимо написать либу, которая будет содержать сервер и клиента для взаимодействия сервисов через unix-socket. 

Требования:
Пакет сервера
Пакет клиента
Readme с описанием и примерами
Размещение на GitHub в публичной репе
Использовать можно только стандартные пакеты go

### Пример использования

#### Запуск сервера
```go
func main() {
	var wg sync.WaitGroup
	cfg := configs.ReadConfig(configPath)

	server := server.NewServer(cfg.ListenAddress) // инициализация сервера
	err := server.Start() // запуск сервера
	if err != nil {
		log.Println(err)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range server.MsgChan { // получение сообщений от сервисов
			log.Printf("%s: %s\n", msg.From, msg.Message)
		}
	}()

        // Cleanup the socket file.
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        wg.Add(1)
		
        go func() { // graceful shutdown
            defer wg.Done()
            <-c
            err := server.Close()
            if err != nil {
                log.Println(err)
                os.Exit(1)
            }
            os.Remove(cfg.ListenAddress)
            os.Exit(0)
        }()
		
        wg.Wait()
}
```

#### Запуск клиента
```go
   func main() {
        client := client2.NewClient(socketPath, "serviceDesk") 
        err := client.ConnectAndWriteToServer()
        if err != nil {
            return
        }
}
```


#### Вывод в консоль
```
2023/06/18 00:36:55 Unix server is ALIVE!
2023/06/18 00:37:37 New connection to the server:  @
2023/06/18 00:37:50 serviceDesk: hello world
2023/06/18 00:39:06 serviceDesk: alert all users that they are super
2023/06/18 00:39:20 New connection to the server:  @
2023/06/18 00:39:47 Telegram Alerts: All users has been notified!
```


