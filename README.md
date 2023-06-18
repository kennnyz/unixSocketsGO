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

### Описание

Реализованы пакеты server и client для взаимодействия через unix сокеты.
#### Структура Сервера:
```go
type Server struct {
    listenAddress string // Путь к сокету
    listener      net.Listener // Слушатель
    msgChan       chan client.Message // Канал куда кидаются все сообщения, которые приходят на сервер
}
```

Для сервера я реализовал 2 основные функции которые работают с клиентами это - acceptLoop и readLoop.

#### acceptLoop
функция которая принимает новые подключения и создает для них отдельные горутины readLoop, которые будут читать сообщения из сокета.

#### readLoop 
функция которая читает сообщения из отдельного подключения, дает response клиенту и отправляет его сообщение в канал msgChan.

-----------------------------------

#### Структура Клиента:
```go
type Client struct {
	ServiceName string
	socketPath  string
	conn        net.Conn
}
```

Для клиента реализована функция WriteToServer, которая отправляет сообщение на сервер и возвращает ответ от сервера.




### Пример использования

#### Запуск сервера
```go
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
```

#### Запуск клиента
```go
    func main() {
        client, err := client2.NewClient(socketPath, "Telegram Alerts")
        if err != nil {
            return
        }
		
        
        reader := bufio.NewReader(os.Stdin)
        for {
            // Запрос ввода отправителя и сообщения
            fmt.Print("Enter your message: ")
            message, _ := reader.ReadString('\n')
            
            // Создание экземпляра структуры Message
            msg := client2.Message{
                From:    strings.TrimSpace(client.ServiceName),
                Message: strings.TrimSpace(message),
            }
            
            err, resp := client.WriteToServer(msg)
            if err != nil {
                log.Println(err)
            }
            
            // Вывод ответа от сервера
            fmt.Println("Server response:", resp)
            
            // Проверка условия завершения
            if strings.TrimSpace(message) == "exit" {
                break
            }
        }
		
    err = client.Close()
    if err != nil {
        log.Println(err)
    }
}
```


#### Вывод в консоль
```
2023/06/18 00:36:55 Unix server is ALIVE!
2023/06/18 00:37:37 New connection to the server:  @ // Подключение клиента
2023/06/18 00:37:50 serviceDesk: hello world // Отправка сообщения сервиса serviceDesk
2023/06/18 00:39:06 serviceDesk: alert all users that they are super
2023/06/18 00:39:20 New connection to the server:  @ // Подключение клиента
2023/06/18 00:39:47 Telegram Alerts: All users has been notified! // Отправка сообщения сервиса Telegram Alerts
```


