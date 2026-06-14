# me

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Lit](https://img.shields.io/badge/Lit-3-324FFF?logo=lit&logoColor=white)](https://lit.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)](https://vitejs.dev/)
[![pnpm](https://img.shields.io/badge/pnpm-10-F69220?logo=pnpm&logoColor=white)](https://pnpm.io/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

## Features

- **Public Site**: プロフィールサイト
- **Admin Dashboard**: 管理画面
- **API Server**: バックエンドAPI

## Tech Stack

| Category              | Technology              |
| --------------------- | ----------------------- |
| **Frontend (Public)** | Lit, Vite, TypeScript   |
| **Frontend (Admin)**  | React, Vite, TypeScript |
| **Backend (API)**     | Go, Air (Live reload)   |
| **Infrastructure**    | Docker, Docker Compose  |
| **Package Manager**   | pnpm (Workspace)        |

## Local Dev

### 1. コンテナのビルドと起動

```sh
docker compose build
docker compose up -d
```

### 2. バックエンド (Go API) の起動

Air を使用してホットリロード付きでAPIサーバーを起動。

```sh
air --proxy.proxy_port "${API_PORT}" -c ./backend/.air.toml | jq .
```

> [!NOTE]
> 別ドメイン構成の動作確認をローカルで行う場合は、バックエンド起動時に `CORS_ALLOWED_ORIGINS` へ許可するフロントエンドのオリジンをカンマ区切りで渡す。
> CORS だけでなく、unsafe method に対するアプリケーション側の Origin 検証にも使用。
>
> ```sh
> CORS_ALLOWED_ORIGINS="http://localhost:5173" air --proxy.proxy_port "${API_PORT}" -c ./backend/.air.toml | jq .
> ```

### 3. フロントエンドの開発サーバー起動

```sh
pnpm install
pnpm dev:frontend
```

## Deployment

- **Frontend**: CloudFront + S3
- **Backend**: APIGateway + Lambda

## Author

[GitHub](https://github.com/umekikazuya)
