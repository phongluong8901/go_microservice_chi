package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	// 1. Tạo dữ liệu mẫu để trả về
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	// 2. Chuyển struct thành chuỗi JSON đẹp (cóụt tab căn lề)
	out, _ := json.MarshalIndent(payload, "", "\t")

	// 3. Thiết lập Header báo cho client biết đây là dữ liệu JSON
	w.Header().Set("Content-Type", "application/json")
	// 4. Trả về mã trạng thái HTTP 202 Accepted
	w.WriteHeader(http.StatusAccepted)

	// 5. Ghi chuỗi JSON ra phản hồi để gửi về cho client
	w.Write(out)

}
