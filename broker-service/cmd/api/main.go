package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

// Bằng cách đưa chúng vào chung một struct, mọi hàm gắn với Config (thông qua (app *Config)) đều có thể dễ dàng truy cập và sử dụng lại các kết nối này mà không cần dùng đến biến toàn cục (global variables).
// 1. Tạo cái "kho" chứa mọi thứ dùng chung
type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ~!")

	app := Config{
		Rabbit: rabbitConn,
	} // Khởi tạo đối tượng

	log.Printf("Starting broker service on port %s", webPort)

	//define http server
	// Có dấu & (&http.Server{...}): Nó tạo ra struct đó trên bộ nhớ và trả về con trỏ (pointer) trỏ đến ô nhớ chứa struct đó (kiểu dữ liệu lúc này sẽ là *http.Server thay vì http.Server).
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(), // Gọi hàm thông qua đối tượng app
	}

	//Start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		// Sửa "amqpp://" thành "amqp://"
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready ...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
	}

	return connection, nil
}
