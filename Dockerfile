# Menggunakan image Golang versi 1.20
FROM golang:1.20-alpine as builder

# Set working directory di dalam container
WORKDIR /app

# Menyalin file go.mod dan go.sum untuk install dependencies
COPY go.mod go.sum ./

# Menjalankan perintah untuk mendownload dependencies
RUN go mod tidy

# Menyalin seluruh kode aplikasi ke dalam container
COPY . .

# Meng-compile aplikasi Golang
RUN go build -o main .

# Menggunakan image Alpine untuk menjalankan aplikasi
FROM alpine:latest

# Install CA certificates (untuk keperluan HTTPS)
RUN apk --no-cache add ca-certificates

# Menyalin binary dari builder ke dalam container
COPY --from=builder /app/main /main

# Menyalin file .env
COPY .env .env

# Expose port yang akan digunakan oleh aplikasi
EXPOSE 8111

# Menjalankan aplikasi
CMD ["/main"]
