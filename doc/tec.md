# BudgetBook Go版 技術スタック

## 目次

1. [プロジェクト概要](#プロジェクト概要)
2. [技術スタック一覧](#技術スタック一覧)
3. [Go](#go)
4. [Gin（Webフレームワーク）](#ginwebフレームワーク)
5. [sqlc（SQLライブラリ）](#sqlcsqlライブラリ)
6. [golang-migrate（マイグレーション）](#golang-migrateマイグレーション)
7. [PostgreSQL](#postgresql)
8. [JWT認証](#jwt認証)
9. [Docker / Docker Compose](#docker--docker-compose)
10. [アーキテクチャ](#アーキテクチャ)
11. [環境セットアップ](#環境セットアップ)

---

## プロジェクト概要

家計簿アプリ「BudgetBook」のGoバックエンド。
オニオンアーキテクチャに準拠し、Domain / Application / Infrastructure / Presentation の4層構成で開発する。

---

## 技術スタック一覧

| カテゴリ | 採用技術 | バージョン |
|----------|----------|-----------|
| 言語 | Go | 1.26.2 |
| Webフレームワーク | Gin | v1.12.0 |
| SQLライブラリ | sqlc | v1.31.0 |
| マイグレーション | golang-migrate | v4.19.1 |
| データベース | PostgreSQL | 17 |
| DBドライバ | pgx | v5 |
| 認証 | JWT + Google OAuth2 | - |
| コンテナ | Docker / Docker Compose | 29.2.1 |

---

## Go

### 概要

Googleが開発したオープンソースのプログラミング言語。シンプルな構文、高速なコンパイル、強力な並行処理が特徴。

### 特徴

- **静的型付け**: コンパイル時に型エラーを検出できる
- **高速**: コンパイル言語のため実行速度が速い
- **並行処理**: `goroutine` と `channel` による軽量な並行処理
- **シンプルな構文**: 予約語が少なく、読みやすいコード
- **標準ライブラリが充実**: HTTPサーバーすら標準ライブラリで実装可能

### 他言語との比較

| 項目 | Go | Java | Python |
|------|-----|------|--------|
| 型付け | 静的 | 静的 | 動的 |
| 実行速度 | 速い | 速い（JVM） | 遅め |
| 並行処理 | goroutine（軽量） | Thread（重め） | asyncio |
| 学習コスト | 低〜中 | 高め | 低い |
| 用途 | API・CLI・インフラ | エンタープライズ | 機械学習・スクリプト |

### Go Modules

依存関係の管理システム。`go.mod` にモジュール名とGoバージョン、依存パッケージを記録する。

```bash
# モジュール初期化
go mod init budget-book-go

# パッケージ追加
go get github.com/gin-gonic/gin

# 不要なパッケージを削除
go mod tidy
```

---

## Gin（Webフレームワーク）

### 概要

GoのWebフレームワーク。GitHubスター数はGoのWebフレームワーク中トップクラスで、国内外問わず最も採用実績が多い。

### 特徴

- **高速**: httprouterベースの高速ルーティング
- **ミドルウェア**: 認証・ロギング・CORSなどを簡単に差し込める
- **シンプルなAPI**: Javaに比べてコード量が少ない
- **情報量が豊富**: 日本語の記事・書籍・チュートリアルが多い

### EchoとGinの比較

| 項目 | Gin | Echo |
|------|-----|------|
| GitHubスター数 | 約80k | 約30k |
| 情報量 | 多い | やや少ない |
| パフォーマンス | 高速 | 同等 |
| ミドルウェア | 豊富 | 豊富 |
| 日本での採用 | 多い | スタートアップに多い |
| 初学者向け | ◎ | ○ |

### 基本的な使い方

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // GETエンドポイント
    r.GET("/api/expenses", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "支出一覧"})
    })

    // パスパラメータ
    r.GET("/api/expenses/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{"id": id})
    })

    r.Run(":8080")
}
```

### ミドルウェアの使い方

```go
// 全ルートに適用
r.Use(gin.Logger())
r.Use(gin.Recovery())

// グループに適用（認証など）
api := r.Group("/api")
api.Use(AuthMiddleware())
{
    api.GET("/expenses", expenseHandler.GetAll)
    api.POST("/expenses", expenseHandler.Create)
}
```

---

## sqlc（SQLライブラリ）

### 概要

SQLを直接書くと、型安全なGoコードを自動生成してくれるツール。

### 特徴

- **型安全**: 生成されたコードは完全に型付けされている
- **SQLそのまま**: ORMの独自構文を覚える必要がない
- **生産性**: クエリを書くだけでGoの関数が生成される
- **N+1対策がしやすい**: JOINを使ったクエリを書くだけ

### GORMとsqlcの比較

| 項目 | GORM | sqlc |
|------|------|------|
| 学習コスト | 低い | 中程度 |
| SQLの知識 | 不要（でも知ってると良い） | 必要 |
| 型安全 | ○ | ◎ |
| パフォーマンス | やや劣る | 高い |
| 複雑なクエリ | 難しい | SQLで自由に書ける |
| 現場採用 | 多い（歴史あり） | 増加中（モダン） |
| 初学者向け | ◎ | ○ |

### 使い方の流れ

```
① query.sql にSQLを書く
      ↓
② sqlc generate を実行
      ↓
③ 型安全なGoコードが自動生成される
      ↓
④ 生成された関数をリポジトリで使う
```

### クエリの書き方

```sql
-- db/queries/expense.sql

-- name: GetExpense :one
SELECT e.*, c.name as category_name
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.id = $1 AND e.user_id = $2;

-- name: ListExpenses :many
SELECT e.*, c.name as category_name
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1
ORDER BY e.expense_date DESC;

-- name: CreateExpense :one
INSERT INTO expenses (user_id, category_id, amount, description, expense_date, payment_method, memo)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
```

---

## golang-migrate（マイグレーション）

### 概要

DBスキーマのバージョン管理ツール。`up` と `down` のSQLファイルをペアで管理し、適用・ロールバックができる。

### 特徴

- **バージョン管理**: どのマイグレーションまで適用済みかDBに記録する
- **ロールバック**: `down` ファイルで元の状態に戻せる
- **CI/CD対応**: コマンド一発でマイグレーション実行できる

### ファイル命名規則

```
db/migrations/
├── 000001_create_users_table.up.sql      # 適用
├── 000001_create_users_table.down.sql    # 巻き戻し
├── 000002_create_categories_table.up.sql
├── 000002_create_categories_table.down.sql
```

### 主要コマンド

```bash
# マイグレーションファイル作成
migrate create -ext sql -dir db/migrations -seq create_users_table

# 全マイグレーションを適用
migrate -path db/migrations -database "postgres://budget_book_user:budget_book_pass@localhost:5432/budget_book_db?sslmode=disable" up

# 1つ前に戻す
migrate -path db/migrations -database "..." down 1

# 現在のバージョン確認
migrate -path db/migrations -database "..." version
```

---

## PostgreSQL

### 概要

オープンソースのRDBMS（リレーショナルデータベース管理システム）。本番環境での採用実績が非常に高く、機能も豊富。

### MySQLとの比較

| 項目 | PostgreSQL | MySQL |
|------|-----------|-------|
| JSON対応 | ◎（jsonb型） | ○ |
| UUID | ◎（ネイティブ対応） | △（文字列で代替） |
| 標準SQL準拠 | 高い | やや低い |
| 日本での採用 | 増加中 | 多い（歴史あり） |
| RailsやDjangoとの相性 | ○ | ○ |

> このプロジェクトではUUID型をPKに使うため、PostgreSQLが最適。

### 接続情報

```yaml
host: localhost
port: 5432
database: budget_book_db
username: budget_book_user
password: budget_book_pass
```

---

## JWT認証

### 概要

JSON Web Token。ユーザー認証情報をトークン（文字列）として発行・検証する仕組み。

### 構造

```
eyJhbGciOiJIUzI1NiJ9  .  eyJ1c2VySWQiOiIxMjMifQ  .  signature
      ヘッダー                    ペイロード               署名
```

### セッション認証との比較

| 項目 | JWT | セッション |
|------|-----|-----------|
| サーバー状態 | 不要（ステートレス） | セッションストアが必要 |
| スケール | 容易 | Redisなどが必要 |
| 無効化 | 難しい（有効期限まで有効） | 即時可能 |
| モバイル対応 | ◎ | △ |

### 使い方（イメージ）

```go
// トークン生成
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "userId": user.ID,
    "exp":    time.Now().Add(24 * time.Hour).Unix(),
})
tokenString, _ := token.SignedString([]byte(secretKey))

// トークン検証（ミドルウェア）
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(secretKey), nil
})
```

### 主要コマンド

```bash
# 起動
docker compose up -d

# 停止
docker compose down

# ログ確認
docker compose logs -f postgres

# DBに接続
docker exec -it budget-book-db psql -U budget_book_user -d budget_book_db
```

---

## アーキテクチャ

### オニオンアーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure                           │
│  (DB接続, pgx/sqlc, 外部API)                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                  Application                        │   │
│  │  (UseCase: ビジネスロジックの流れを制御)              │   │
│  │  ┌─────────────────────────────────────────────┐   │   │
│  │  │                 Domain                      │   │   │
│  │  │  (Entity, Repository Interface, ValueObject)│   │   │
│  │  │                                             │   │   │
│  │  │   ★ここが中心。外部に依存しない             │   │   │
│  │  └─────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘

依存の方向: 外側 → 内側 (Infrastructure → Application → Domain)
```

### ディレクトリ構成

```
budget-book-go/
├── cmd/api/main.go                          # エントリーポイント
├── internal/
│   ├── domain/                              # ドメイン層（最も重要）
│   │   ├── entity/                          # ビジネスエンティティ
│   │   ├── valueobject/                     # 値オブジェクト（Money等）
│   │   ├── repository/                      # リポジトリインターフェース
│   │   └── error/                           # ドメインエラー定義
│   ├── application/                         # アプリケーション層
│   │   ├── usecase/                         # ユースケース
│   │   └── dto/                             # データ転送オブジェクト
│   ├── infrastructure/                      # インフラ層
│   │   ├── persistence/postgres/            # リポジトリ実装
│   │   ├── persistence/sqlc/                # sqlc生成コード
│   │   └── config/                          # 設定読み込み
│   └── presentation/                        # プレゼンテーション層
│       ├── handler/                         # HTTPハンドラ
│       ├── request/                         # リクエスト構造体
│       ├── response/                        # レスポンス構造体
│       └── middleware/                      # 認証ミドルウェア等
└── db/
    └── migrations/                          # マイグレーションSQL
```

---

## 環境セットアップ

### 前提条件

| ツール | バージョン | インストール方法 |
|--------|-----------|----------------|
| Go | 1.26.2 | `brew install go` |
| sqlc | 1.31.0 | `brew install sqlc` |
| golang-migrate | 4.19.1 | `brew install golang-migrate` |
| Docker | 29.2.1 | Docker Desktop |

### セットアップ手順

```bash
# 1. リポジトリのクローン（またはディレクトリ作成）
cd ~/systems/budget-book-go

# 2. Goモジュール初期化
go mod init budget-book-go

# 3. 依存パッケージのインストール
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5
go get github.com/google/uuid
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv

# 4. PostgreSQL起動
docker compose up -d

# 5. マイグレーション実行
migrate -path db/migrations \
  -database "postgres://budget_book_user:budget_book_pass@localhost:5432/budget_book_db?sslmode=disable" \
  up

# 6. sqlcコード生成
sqlc generate

# 7. サーバー起動
go run cmd/api/main.go
```

### Makefile（便利コマンド集）

```makefile
.PHONY: up down migrate sqlc run

up:
	docker compose up -d

down:
	docker compose down

migrate:
	migrate -path db/migrations \
	  -database "postgres://budget_book_user:budget_book_pass@localhost:5432/budget_book_db?sslmode=disable" \
	  up

sqlc:
	sqlc generate

run:
	go run cmd/api/main.go
```