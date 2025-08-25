FROM golang:1.25-alpine

WORKDIR /app

# 必要なパッケージのインストール
RUN apk add --no-cache gcc musl-dev git

# go.modとgo.sumを先にコピーして依存解決
COPY go.mod go.sum ./
RUN go mod tidy

# ソースコード全体をコピー
COPY . .

# ビルド
RUN go build -o main .

# デフォルトポートを8080に変更
EXPOSE 8080

CMD ["./main"] 
