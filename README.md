# Budget Book Go - 家計簿管理 API

個人の収支を管理するための REST API サーバー。Google OAuth2 認証、支出/収入の CRUD、定期支出の自動生成、レシート OCR（Google Gemini）、収支予測などの機能を提供する。

---

## 目次

- [技術スタック](#技術スタック)
- [プロジェクト構成](#プロジェクト構成)
- [アーキテクチャ](#アーキテクチャ)
- [レイヤー詳細](#レイヤー詳細)
- [ER図（データベーススキーマ）](#er図データベーススキーマ)
- [API エンドポイント一覧](#api-エンドポイント一覧)
- [リクエストフロー](#リクエストフロー)
- [認証フロー](#認証フロー)
- [OCR（レシート解析）フロー](#ocrレシート解析フロー)
- [定期支出の自動生成](#定期支出の自動生成)
- [セットアップ](#セットアップ)
- [開発コマンド](#開発コマンド)
- [タスク別ガイド：どこを編集すればいいか](#タスク別ガイドどこを編集すればいいか)

---

## 技術スタック

| カテゴリ | 技術 |
|---|---|
| 言語 | Go 1.26 |
| HTTP フレームワーク | Gin |
| データベース | PostgreSQL 17 |
| DB ドライバ | pgx/v5 (pgxpool) |
| SQL コード生成 | sqlc |
| 認証 | Google OAuth2 + JWT (HMAC-SHA256) |
| AI/OCR | Google Gemini API (`gemini-3.1-flash-lite`) |
| コンテナ | Docker / Docker Compose |
| マイグレーション | golang-migrate |

---

## プロジェクト構成

```
budget-book-go/
├── cmd/api/
│   └── main.go                    # エントリーポイント（DI・ルーティング定義）
│
├── internal/
│   ├── presentation/              # 🌐 プレゼンテーション層（HTTP）
│   │   ├── handler/               #   各リソースの HTTP ハンドラ
│   │   ├── request/               #   リクエスト構造体（バリデーション付き）
│   │   ├── response/              #   レスポンス構造体
│   │   └── middleware/            #   JWT 認証ミドルウェア
│   │
│   ├── application/               # 📋 アプリケーション層（ユースケース）
│   │   ├── usecase/               #   ビジネスロジックの実行単位
│   │   │   ├── expense/           #     支出の作成・取得・更新・削除
│   │   │   ├── income/            #     収入の作成・取得・更新・削除
│   │   │   ├── category/          #     カテゴリ管理
│   │   │   ├── recurring_expense/ #     定期支出の管理・自動生成
│   │   │   ├── summary/           #     収支予測（フォーキャスト）
│   │   │   ├── auth/              #     Google OAuth 認証
│   │   │   └── ocr/               #     レシート OCR 解析
│   │   └── dto/                   #   データ転送オブジェクト（Command/Result）
│   │
│   ├── domain/                    # 🏛️ ドメイン層（ビジネスルール）
│   │   ├── entity/                #   ドメインエンティティ
│   │   ├── repository/            #   リポジトリインターフェース（契約）
│   │   ├── valueobject/           #   値オブジェクト（Money）
│   │   └── error/                 #   ドメインエラー定義
│   │
│   └── infrastructure/            # 🔧 インフラ層（外部システム連携）
│       ├── config/                #   設定（DB, OAuth, JWT, Gemini）
│       └── persistence/           #   データベースアクセス
│           ├── postgres/          #     リポジトリ実装（pgx）
│           └── sqlc/              #     sqlc 自動生成コード
│
├── db/
│   ├── migrations/                # マイグレーションファイル
│   ├── queries/                   # sqlc 用 SQL クエリ定義
│   ├── schema/                    # スキーマ定義
│   └── seed.sql                   # シードデータ
│
├── Dockerfile                     # マルチステージビルド
├── docker-compose.yml             # PostgreSQL コンテナ
├── Makefile                       # 開発コマンド集
├── sqlc.yaml                      # sqlc 設定
└── .env.example                   # 環境変数テンプレート
```

---

## アーキテクチャ

クリーンアーキテクチャ（Hexagonal Architecture）を採用。依存の方向は外側から内側へ一方向。

```
┌─────────────────────────────────────────────────────────────────────┐
│                        External Systems                             │
│              (PostgreSQL, Google OAuth, Gemini API)                  │
└───────────────────────────────┬─────────────────────────────────────┘
                                │
┌───────────────────────────────▼─────────────────────────────────────┐
│                     Infrastructure Layer                             │
│                                                                      │
│  config/          persistence/postgres/       persistence/sqlc/      │
│  (DB,OAuth,JWT,   (リポジトリ実装)             (自動生成コード)       │
│   Gemini設定)                                                        │
│                                                                      │
│  ※ domain/repository のインターフェースを実装                        │
└───────────────────────────────┬─────────────────────────────────────┘
                                │ implements
┌───────────────────────────────▼─────────────────────────────────────┐
│                        Domain Layer                                  │
│                                                                      │
│  entity/              repository/          valueobject/    error/     │
│  (User, Expense,      (インターフェース     (Money)        (Domain    │
│   Income, Category,    定義のみ)                            Error)    │
│   RecurringExpense)                                                   │
│                                                                      │
│  ※ 外部依存なし。純粋なビジネスルール                                │
└───────────────────────────────▲─────────────────────────────────────┘
                                │ depends on
┌───────────────────────────────┴─────────────────────────────────────┐
│                      Application Layer                               │
│                                                                      │
│  usecase/                              dto/                          │
│  (CreateExpense, GetExpense,           (Command: 入力用DTO           │
│   GoogleAuth, AnalyzeReceipt,           Result: 出力用DTO)           │
│   GetForecast, ApplyRecurring ...)                                   │
│                                                                      │
│  ※ ユースケースごとに1ファイル。リポジトリIF経由でデータ操作          │
└───────────────────────────────▲─────────────────────────────────────┘
                                │ calls
┌───────────────────────────────┴─────────────────────────────────────┐
│                     Presentation Layer                                │
│                                                                      │
│  handler/             request/           response/     middleware/    │
│  (ExpenseHandler,     (CreateExpense     (JSON応答     (JWT認証       │
│   IncomeHandler,       Request ...)      構造体)       ミドルウェア)  │
│   AuthHandler ...)                                                    │
│                                                                      │
│  ※ HTTP リクエスト/レスポンスの変換。Gin に依存                      │
└───────────────────────────────▲─────────────────────────────────────┘
                                │ HTTP
┌───────────────────────────────┴─────────────────────────────────────┐
│                          Client                                      │
│                    (Frontend / curl / etc.)                           │
└─────────────────────────────────────────────────────────────────────┘
```

**依存の方向:**
```
Client → Presentation → Application → Domain ← Infrastructure
```

> Domain 層はどこにも依存しない。Infrastructure 層が Domain 層のインターフェースを実装する（依存性逆転の原則）。

---

## レイヤー詳細

### Presentation 層 (`internal/presentation/`)

HTTP リクエストの受け付けとレスポンス返却を担当。

| ディレクトリ | 役割 | 主なファイル |
|---|---|---|
| `handler/` | エンドポイントの処理 | `expense_handler.go`, `income_handler.go`, `auth_handler.go`, `ocr_handler.go` 等 |
| `request/` | リクエストボディのバインドとバリデーション | `expense_request.go`, `income_request.go` 等 |
| `response/` | レスポンス構造体 | レスポンスの JSON 構造定義 |
| `middleware/` | 認証ミドルウェア | `auth.go` — JWT 検証 & ユーザー ID 抽出 |

### Application 層 (`internal/application/`)

ビジネスロジックのオーケストレーション。ユースケースごとに1構造体。

| ユースケース | ファイル | 概要 |
|---|---|---|
| 支出作成 | `usecase/expense/create_expense.go` | Money バリデーション → Entity 生成 → Repository 保存 |
| 支出取得 | `usecase/expense/get_expense.go` | 一覧/ID指定/日付範囲/検索 の各メソッド |
| 支出更新 | `usecase/expense/update_expense.go` | 既存支出の更新 |
| 支出削除 | `usecase/expense/delete_expense.go` | 支出の削除 |
| 収入 CRUD | `usecase/income/` | 支出と同様の構造 |
| カテゴリ CRUD | `usecase/category/` | カテゴリの管理 |
| 定期支出管理 | `usecase/recurring_expense/` | CRUD + 自動生成ロジック |
| 収支予測 | `usecase/summary/get_forecast.go` | 現在残高・予測残高の計算 |
| 認証 | `usecase/auth/google_auth.go` | OAuth2 フロー + JWT 発行 |
| OCR | `usecase/ocr/analyze_receipt.go` | Gemini API でレシート解析 |

### Domain 層 (`internal/domain/`)

純粋なビジネスルール。外部ライブラリへの依存なし。

| 要素 | ファイル | 内容 |
|---|---|---|
| User | `entity/user.go` | ID, Email, Provider, Name 等 |
| Expense | `entity/expense.go` | 支出（金額, 日付, カテゴリ, 支払方法） |
| Income | `entity/income.go` | 収入 |
| Category | `entity/category.go` | カテゴリ（EXPENSE/INCOME の type 区分） |
| RecurringExpense | `entity/recurring_expense.go` | 定期支出（BillingDay, StartDate, EndDate, IsActive） |
| Money | `valueobject/money.go` | 金額の値オブジェクト（0より大きいことを保証） |
| Repository IF | `repository/*.go` | 各エンティティのリポジトリインターフェース |
| DomainError | `error/domain_error.go` | NOT_FOUND / UNAUTHORIZED / INVALID_INPUT |

### Infrastructure 層 (`internal/infrastructure/`)

外部システムとの接続。

| ディレクトリ | 役割 |
|---|---|
| `config/` | DB接続 (`config.go`), OAuth設定 (`oauth.go`), JWT秘密鍵 (`jwt.go`), Gemini APIキー (`gemini.go`) |
| `persistence/postgres/` | リポジトリインターフェースの PostgreSQL 実装 |
| `persistence/sqlc/` | sqlc による自動生成コード（クエリ実行関数、パラメータ/結果型） |

---

## ER図（データベーススキーマ）

```
┌──────────────────────┐
│        users          │
├──────────────────────┤
│ id          UUID  PK  │
│ email       VARCHAR   │◄─────────────────────────────────────┐
│ password_hash VARCHAR │                                      │
│ provider    VARCHAR   │  ┌────────────────────────┐          │
│ provider_id VARCHAR   │  │     categories          │          │
│ name        VARCHAR   │  ├────────────────────────┤          │
│ created_at  TIMESTAMP │  │ id         UUID  PK     │          │
│ updated_at  TIMESTAMP │  │ user_id    UUID  FK ────┤──────────┤
└──────────────────────┘  │ name       VARCHAR      │          │
                           │ type       VARCHAR      │ EXPENSE  │
                           │            (CHECK)      │ / INCOME │
                           │ color      VARCHAR      │          │
                           │ sort_order INTEGER      │          │
                           │ is_default BOOLEAN      │          │
                           │ created_at TIMESTAMP    │          │
                           └─────────┬──────────────┘          │
                                     │                          │
                    ┌────────────────┼────────────────┐         │
                    │                │                │         │
        ┌───────────▼────────┐ ┌─────▼──────────┐ ┌──▼─────────────────────┐
        │     expenses        │ │    incomes      │ │  recurring_expenses     │
        ├────────────────────┤ ├────────────────┤ ├────────────────────────┤
        │ id        UUID  PK  │ │ id     UUID PK  │ │ id          UUID  PK   │
        │ user_id   UUID  FK ─┤►│ user_id UUID FK─┤►│ user_id     UUID  FK ──┤►users
        │ category_id UUID FK │ │ category_id FK  │ │ category_id UUID  FK   │
        │ amount    DECIMAL   │ │ amount  DECIMAL │ │ amount      DECIMAL    │
        │ description VARCHAR │ │ description     │ │ description VARCHAR    │
        │ expense_date DATE   │ │ income_date     │ │ payment_method VARCHAR │
        │ payment_method      │ │ memo    TEXT    │ │ memo        TEXT       │
        │   VARCHAR           │ │ created_at      │ │ billing_day INTEGER    │
        │ memo      TEXT      │ │ updated_at      │ │   (1-31)               │
        │ created_at          │ └────────────────┘ │ start_date  DATE       │
        │ updated_at          │                     │ end_date    DATE       │
        └──────────┬─────────┘                     │ is_active   BOOLEAN    │
                   │                                │ created_at             │
                   │                                │ updated_at             │
                   │                                └───────────┬───────────┘
                   │                                            │
                   │    ┌──────────────────────────────┐        │
                   │    │  recurring_expense_logs        │        │
                   │    ├──────────────────────────────┤        │
                   │    │ id                    UUID PK │        │
                   └────┤ expense_id            UUID FK │        │
                        │ recurring_expense_id  UUID FK ┤────────┘
                        │ billing_year          INTEGER │
                        │ billing_month         INTEGER │
                        │ created_at            TIMESTAMP│
                        │                                │
                        │ UNIQUE(recurring_expense_id,   │
                        │        billing_year,           │
                        │        billing_month)          │
                        └───────────────────────────────┘
```

**主な制約:**
- `expenses.amount > 0`, `incomes.amount > 0`, `recurring_expenses.amount > 0`
- `expenses.payment_method IN ('CASH', 'CREDIT_CARD')`
- `categories.type IN ('EXPENSE', 'INCOME')`
- `recurring_expenses.billing_day BETWEEN 1 AND 31`
- `categories` は `(user_id, name, type)` でユニーク

---

## API エンドポイント一覧

ベース URL: `/api`

### 認証（公開）

| メソッド | パス | 概要 |
|---|---|---|
| `GET` | `/auth/google/url` | Google OAuth URL を取得 |
| `GET` | `/auth/google/callback` | OAuth コールバック（JWT 発行 → フロントにリダイレクト） |

### 支出（要認証）

| メソッド | パス | 概要 |
|---|---|---|
| `GET` | `/expenses` | ユーザーの全支出を取得（定期支出の自動生成も実行） |
| `GET` | `/expenses/:id` | 支出を ID で取得 |
| `GET` | `/expenses/planned` | 未来日付の予定支出を取得 |
| `GET` | `/expenses/date?from=&to=` | 日付範囲で支出を取得 |
| `GET` | `/expenses/search?from=&to=&categoryId=&keyword=` | 条件で支出を検索 |
| `POST` | `/expenses` | 支出を作成 |
| `PUT` | `/expenses/:id` | 支出を更新 |
| `DELETE` | `/expenses/:id` | 支出を削除 |

### 収入（要認証）

| メソッド | パス | 概要 |
|---|---|---|
| `GET` | `/incomes` | ユーザーの全収入を取得 |
| `GET` | `/incomes/:id` | 収入を ID で取得 |
| `GET` | `/incomes/planned` | 未来日付の予定収入を取得 |
| `GET` | `/incomes/date?from=&to=` | 日付範囲で収入を取得 |
| `POST` | `/incomes` | 収入を作成 |
| `PUT` | `/incomes/:id` | 収入を更新 |
| `DELETE` | `/incomes/:id` | 収入を削除 |

### カテゴリ（要認証）

| メソッド | パス | 概要 |
|---|---|---|
| `GET` | `/categories` | ユーザーのカテゴリ一覧 |
| `POST` | `/categories` | カテゴリ作成 |
| `PUT` | `/categories/:id` | カテゴリ更新 |
| `DELETE` | `/categories/:id` | カテゴリ削除 |

### 定期支出（要認証）

| メソッド | パス | 概要 |
|---|---|---|
| `GET` | `/recurring-expenses` | 定期支出一覧 |
| `GET` | `/recurring-expenses/:id` | 定期支出を ID で取得 |
| `POST` | `/recurring-expenses` | 定期支出を作成 |
| `PUT` | `/recurring-expenses/:id` | 定期支出を更新 |
| `DELETE` | `/recurring-expenses/:id` | 定期支出を削除 |

### サマリー（要認証）

| メソッド | パス | 概要 |
|---|---|---|
| `GET` | `/summary/forecast?targetDate=` | 収支予測を取得 |

### OCR（要認証）

| メソッド | パス | 概要 |
|---|---|---|
| `POST` | `/ocr/analyze` | レシート画像を解析（multipart form-data） |

---

## リクエストフロー

### 支出作成の例: `POST /api/expenses`

```
Client                Middleware          Handler              UseCase             Repository          Database
  │                      │                  │                    │                    │                  │
  │  POST /api/expenses  │                  │                    │                    │                  │
  │  Authorization:      │                  │                    │                    │                  │
  │  Bearer <JWT>        │                  │                    │                    │                  │
  │─────────────────────►│                  │                    │                    │                  │
  │                      │                  │                    │                    │                  │
  │                      │ JWT検証           │                    │                    │                  │
  │                      │ userID抽出        │                    │                    │                  │
  │                      │ contextに保存     │                    │                    │                  │
  │                      │─────────────────►│                    │                    │                  │
  │                      │                  │                    │                    │                  │
  │                      │                  │ JSONパース          │                    │                  │
  │                      │                  │ バリデーション       │                    │                  │
  │                      │                  │ Command DTO生成     │                    │                  │
  │                      │                  │───────────────────►│                    │                  │
  │                      │                  │                    │                    │                  │
  │                      │                  │                    │ Money検証           │                  │
  │                      │                  │                    │ (amount > 0)       │                  │
  │                      │                  │                    │ Entity生成          │                  │
  │                      │                  │                    │───────────────────►│                  │
  │                      │                  │                    │                    │                  │
  │                      │                  │                    │                    │ Go型→pgtype変換   │
  │                      │                  │                    │                    │ INSERT実行         │
  │                      │                  │                    │                    │─────────────────►│
  │                      │                  │                    │                    │                  │
  │                      │                  │                    │                    │◄─────────────────│
  │                      │                  │                    │                    │ 結果→Entity変換   │
  │                      │                  │                    │◄───────────────────│                  │
  │                      │                  │                    │                    │                  │
  │                      │                  │                    │ Entity→Result DTO  │                  │
  │                      │                  │◄───────────────────│                    │                  │
  │                      │                  │                    │                    │                  │
  │                      │                  │ JSON レスポンス     │                    │                  │
  │◄────────────────────────────────────────│ 201 Created       │                    │                  │
  │                      │                  │                    │                    │                  │
```

**対応ファイル:**

| ステップ | ファイル |
|---|---|
| JWT 検証 | `internal/presentation/middleware/auth.go` |
| リクエスト解析 | `internal/presentation/handler/expense_handler.go` → `Create()` |
| リクエスト型 | `internal/presentation/request/expense_request.go` |
| DTO 変換 | `internal/application/dto/expense_dto.go` |
| ビジネスロジック | `internal/application/usecase/expense/create_expense.go` |
| 値バリデーション | `internal/domain/valueobject/money.go` |
| DB 操作 | `internal/infrastructure/persistence/postgres/expense_repository.go` → `Save()` |
| SQL クエリ | `db/queries/expense.sql` → `CreateExpense` |

---

## 認証フロー

Google OAuth2 + JWT によるステートレス認証。

```
     Client                  API Server              Google               Database
       │                        │                      │                     │
  ①    │  GET /auth/google/url  │                      │                     │
       │──────────────────────►│                      │                     │
       │                        │ state生成             │                     │
       │◄──────────────────────│ OAuth URL返却         │                     │
       │  { "url": "https://   │                      │                     │
       │    accounts.google..." │                      │                     │
       │  }                     │                      │                     │
       │                        │                      │                     │
  ②    │  ブラウザでGoogleログイン │                      │                     │
       │───────────────────────────────────────────────►│                     │
       │                        │                      │                     │
  ③    │  リダイレクト           │                      │                     │
       │  /auth/google/callback │                      │                     │
       │  ?code=xxx&state=xxx   │                      │                     │
       │──────────────────────►│                      │                     │
       │                        │                      │                     │
  ④    │                        │ code→token交換        │                     │
       │                        │─────────────────────►│                     │
       │                        │◄─────────────────────│                     │
       │                        │ access_token取得      │                     │
       │                        │                      │                     │
  ⑤    │                        │ ユーザー情報取得       │                     │
       │                        │─────────────────────►│                     │
       │                        │◄─────────────────────│                     │
       │                        │ email, name, id      │                     │
       │                        │                      │                     │
  ⑥    │                        │ ユーザー検索/作成      │                     │
       │                        │──────────────────────────────────────────►│
       │                        │◄──────────────────────────────────────────│
       │                        │                      │                     │
  ⑦    │                        │ JWT生成               │                     │
       │                        │ (userId, email,      │                     │
       │                        │  有効期限: 24h)       │                     │
       │                        │                      │                     │
  ⑧    │◄──────────────────────│                      │                     │
       │  リダイレクト:          │                      │                     │
       │  {FRONTEND_URL}/       │                      │                     │
       │  auth/callback?        │                      │                     │
       │  token=xxx&            │                      │                     │
       │  userId=xxx&           │                      │                     │
       │  email=xxx&            │                      │                     │
       │  name=xxx              │                      │                     │
       │                        │                      │                     │
  以降  │  Authorization:        │                      │                     │
       │  Bearer <JWT>          │                      │                     │
       │──────────────────────►│ middleware で検証      │                     │
       │                        │                      │                     │
```

**対応ファイル:**

| ステップ | ファイル |
|---|---|
| OAuth URL 生成 | `internal/presentation/handler/auth_handler.go` → `GetGoogleAuthURL()` |
| OAuth 設定 | `internal/infrastructure/config/oauth.go` |
| コールバック処理 | `internal/presentation/handler/auth_handler.go` → `GoogleCallback()` |
| Token交換 & ユーザー作成 | `internal/application/usecase/auth/google_auth.go` |
| JWT 生成 | `internal/application/usecase/auth/google_auth.go` |
| JWT 設定 | `internal/infrastructure/config/jwt.go` |
| JWT 検証 | `internal/presentation/middleware/auth.go` |

**Mock モード:** `IS_MOCK=True` 設定時は JWT を検証せず、`X-User-Id` ヘッダーまたは `MOCK_USER_ID` 環境変数からユーザー ID を取得する（開発・テスト用）。

---

## OCR（レシート解析）フロー

```
     Client                    Handler               UseCase             Gemini API
       │                         │                      │                    │
       │  POST /ocr/analyze      │                      │                    │
       │  Content-Type:          │                      │                    │
       │  multipart/form-data    │                      │                    │
       │  [image file]           │                      │                    │
       │────────────────────────►│                      │                    │
       │                         │                      │                    │
       │                         │ 画像データ抽出        │                    │
       │                         │ MIMEタイプ判定        │                    │
       │                         │─────────────────────►│                    │
       │                         │                      │                    │
       │                         │                      │ プロンプト + 画像   │
       │                         │                      │ → Gemini API       │
       │                         │                      │───────────────────►│
       │                         │                      │                    │
       │                         │                      │◄───────────────────│
       │                         │                      │ JSON応答            │
       │                         │                      │                    │
       │                         │                      │ JSONパース          │
       │                         │                      │ OCRResult生成      │
       │                         │◄─────────────────────│                    │
       │                         │                      │                    │
       │◄────────────────────────│                      │                    │
       │  200 OK                 │                      │                    │
       │  {                      │                      │                    │
       │    "items": [           │                      │                    │
       │      {                  │                      │                    │
       │        "description":   │                      │                    │
       │          "コーヒー",     │                      │                    │
       │        "amount": 350,   │                      │                    │
       │        "expenseDate":   │                      │                    │
       │          "2026-10-17",  │                      │                    │
       │        "paymentMethod": │                      │                    │
       │          "CASH"         │                      │                    │
       │      }                  │                      │                    │
       │    ]                    │                      │                    │
       │  }                      │                      │                    │
```

**対応ファイル:**

| ステップ | ファイル |
|---|---|
| 画像受け取り | `internal/presentation/handler/ocr_handler.go` |
| Gemini 呼び出し | `internal/application/usecase/ocr/analyze_receipt.go` |
| API キー設定 | `internal/infrastructure/config/gemini.go` |

---

## 定期支出の自動生成

定期支出（サブスク等）は `GET /api/expenses` のタイミングでオンデマンドに自動生成される。

```
┌─ GET /api/expenses ─────────────────────────────────────────────┐
│                                                                  │
│  1. ApplyRecurringExpenseUseCase.Execute()                       │
│     │                                                            │
│     ├─ アクティブな定期支出を全件取得                              │
│     │                                                            │
│     └─ 各定期支出について:                                        │
│        │                                                         │
│        ├─ start_date ～ 現在月まで繰り返し                        │
│        │   │                                                     │
│        │   ├─ recurring_expense_logs に当月レコードあり？          │
│        │   │   ├─ YES → スキップ（重複防止）                      │
│        │   │   └─ NO  → ↓                                       │
│        │   │                                                     │
│        │   ├─ billing_day が月末を超える場合は月末日に調整         │
│        │   │   (例: billing_day=31, 2月 → 2月28日)               │
│        │   │                                                     │
│        │   ├─ expenses テーブルに支出レコードを INSERT             │
│        │   │                                                     │
│        │   └─ recurring_expense_logs に記録                       │
│        │       (recurring_expense_id, billing_year, billing_month)│
│        │                                                         │
│        └─ end_date を超えたら終了                                 │
│                                                                  │
│  2. GetExpenseUseCase.ExecuteGetByUserID()                       │
│     └─ 通常の支出一覧取得（自動生成分も含む）                     │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

**対応ファイル:** `internal/application/usecase/recurring_expense/apply_recurring_expense.go`

---

## セットアップ

### 1. 環境変数の準備

```bash
cp .env.example .env
# .env を編集して各値を設定
```

### 2. PostgreSQL 起動

```bash
make up        # docker compose up -d
```

### 3. マイグレーション実行

```bash
make migrate   # DB スキーマを作成
```

### 4. （任意）シードデータ投入

```bash
make seed
```

### 5. サーバー起動

```bash
make run       # go run cmd/api/main.go
```

サーバーは `PORT` 環境変数のポート（デフォルト: `8080`）で起動する。

---

## 開発コマンド

```bash
make up              # PostgreSQL コンテナ起動
make down            # PostgreSQL コンテナ停止
make migrate         # マイグレーション実行
make migrate-down    # マイグレーション1つロールバック
make seed            # シードデータ投入
make sqlc            # SQL クエリから Go コードを自動生成
make run             # サーバー起動
```

---

## タスク別ガイド：どこを編集すればいいか

### 新しい API エンドポイントを追加したい

1. **SQL クエリ定義** — `db/queries/<resource>.sql` にクエリを追加
2. **sqlc 再生成** — `make sqlc` を実行
3. **リポジトリ IF** — `internal/domain/repository/<resource>_repository.go` にメソッド追加
4. **リポジトリ実装** — `internal/infrastructure/persistence/postgres/<resource>_repository.go` に実装追加
5. **DTO** — `internal/application/dto/<resource>_dto.go` に Command/Result 追加
6. **ユースケース** — `internal/application/usecase/<resource>/` に新ファイル作成
7. **リクエスト型** — `internal/presentation/request/<resource>_request.go` に追加
8. **ハンドラ** — `internal/presentation/handler/<resource>_handler.go` にメソッド追加
9. **ルーティング** — `cmd/api/main.go` にルート登録

### 既存 API のビジネスロジックを変更したい

→ `internal/application/usecase/<resource>/` 内の該当ユースケースファイル

### DB スキーマを変更したい

1. `db/migrations/` に新しいマイグレーションファイルを作成
2. `db/schema/schema.sql` も更新
3. 影響する `db/queries/*.sql` を更新
4. `make sqlc` で再生成
5. リポジトリ実装を更新

### バリデーションを追加・変更したい

- **リクエストレベル:** `internal/presentation/request/` の構造体タグ（Gin binding）
- **ドメインレベル:** `internal/domain/valueobject/` または `internal/domain/entity/`
- **DB レベル:** `db/migrations/` に CHECK 制約を追加

### 新しい外部サービスを統合したい

1. `internal/infrastructure/config/` に設定ファイル追加
2. `internal/application/usecase/` にユースケース作成
3. `internal/presentation/handler/` にハンドラ追加
4. `cmd/api/main.go` で DI & ルーティング

### 認証周りを変更したい

- **ミドルウェア:** `internal/presentation/middleware/auth.go`
- **OAuth フロー:** `internal/application/usecase/auth/google_auth.go`
- **JWT 設定:** `internal/infrastructure/config/jwt.go`
- **OAuth 設定:** `internal/infrastructure/config/oauth.go`
