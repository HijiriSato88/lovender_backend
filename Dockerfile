##############################
# Cloud Run用（本番環境）
##############################

####################
# ビルドステージ
####################
FROM golang:1.23-alpine AS builder

WORKDIR /app

# インストール可能なパッケージ一覧の更新
RUN apk update && \
    apk upgrade && \
    # パッケージのインストール（--no-cacheでキャッシュ削除）
    apk add --no-cache \
            git \
            curl

# Go モジュールファイルをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド（静的リンク）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

####################
# 実行ステージ
####################
FROM alpine:latest

WORKDIR /app

# インストール可能なパッケージ一覧の更新
RUN apk update && \
    apk upgrade && \
    # パッケージのインストール（--no-cacheでキャッシュ削除）
    apk add --no-cache \
            ca-certificates \
            tzdata \
            curl

# ビルドステージで作成したバイナリをコピー
COPY --from=builder /app/main .

# Cloud Run用にポートを動的に設定
EXPOSE 8080

# 実行
CMD ["./main"]
