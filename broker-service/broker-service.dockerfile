# Lấy image golang phiên bản 1.18 trên nền alpine và đặt tên giai đoạn là "builder" (thêm r ở đây)
FROM golang:1.25-alpine as builder

# Tạo thư mục /app trong container
RUN mkdir /app

# Copy toàn bộ mã nguồn ở thư mục hiện tại trên máy vào thư mục /app trong container
COPY . /app

# Chuyển hướng thư mục làm việc hiện tại vào /app
WORKDIR /app

# Biên dịch ứng dụng Go thành file chạy có tên là "brokerApp"
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# Phân quyền thực thi cho file vừa build
RUN chmod +x /app/brokerApp

# Lấy image alpine:latest để làm môi trường chạy thực tế
FROM alpine:latest

# Tạo thư mục /app ở image mới này
RUN mkdir /app

# Copy file binary từ giai đoạn "builder" sang đây (đã khớp tên)
COPY --from=builder /app/brokerApp /app

# Lệnh sẽ tự động chạy khi container khởi động
CMD ["/app/brokerApp"]