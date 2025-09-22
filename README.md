# lovender Backend

Docker + Golang + Atlas を使った3層アーキテクチャ構成

## 技術スタック
- **言語**: Go 1.23
- **フレームワーク**: Echo v4
- **データベース**: MySQL 8.0
- **DB接続**: 標準SQLパッケージ
- **スキーマ管理**: Atlas
- **コンテナ**: Docker

## セットアップ

### 1. 初回プロジェクトを起動

```bash
# Dockerコンテナを起動
docker-compose up --build

# バックグラウンドで起動する場合
docker-compose up -d --build
```
初回起動時に以下が自動で実行されます：
- テーブル作成
- サンプルデータの投入

### 2. 初回以降のプロジェクトを起動

```bash
# Dockerコンテナを起動
docker-compose up

# バックグラウンドで起動する場合
docker-compose up -d
```

### 3. 開発環境でアプリを起動

```bash
go run cmd/server/main.go
```

### 4. APIテスト
以下をコピペしてターミナルで叩く
```
curl http://localhost:8080/api/users/1
```
成功レスポンス
```
{"id":1,"name":"田中太郎","email":"tanaka@example.com","created_at":"2025-09-22T12:16:36.117Z","updated_at":"2025-09-22T12:16:36.117Z"}
```

### 環境変数

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | データベースホスト |
| DB_PORT | 3306 | データベースポート |
| DB_USER | lovender_user | データベースユーザー |
| DB_PASSWORD | lovender_password | データベースパスワード |
| DB_NAME | lovender | データベース名 |
