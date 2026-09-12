FROM golang:1.21-alpine AS builder

WORKDIR /app

# Mengunduh dependensi
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh kode sumber
COPY . .

# Membangun aplikasi
RUN go build -o api cmd/api/main.go

# Image final yang ringan
FROM alpine:latest
WORKDIR /app

# Menyalin file binary dari tahap builder
COPY --from=builder /app/api .
COPY --from=builder /app/docs ./docs

# Azure mengharapkan port 8080 secara default untuk Web App
EXPOSE 8080

# Menjalankan aplikasi
CMD ["./api"]
