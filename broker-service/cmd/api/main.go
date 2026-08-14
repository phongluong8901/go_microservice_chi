package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

// Bằng cách đưa chúng vào chung một struct, mọi hàm gắn với Config (thông qua (app *Config)) đều có thể dễ dàng truy cập và sử dụng lại các kết nối này mà không cần dùng đến biến toàn cục (global variables).
// 1. Tạo cái "kho" chứa mọi thứ dùng chung
type Config struct{}

func main() {
	app := Config{} // Khởi tạo đối tượng

	log.Printf("Starting broker service on port %s", webPort)

	//define http server
	// Có dấu & (&http.Server{...}): Nó tạo ra struct đó trên bộ nhớ và trả về con trỏ (pointer) trỏ đến ô nhớ chứa struct đó (kiểu dữ liệu lúc này sẽ là *http.Server thay vì http.Server).
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(), // Gọi hàm thông qua đối tượng app
	}

	//Start thr server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
