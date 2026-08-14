# me

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Lit](https://img.shields.io/badge/Lit-3-324FFF?logo=lit&logoColor=white)](https://lit.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)](https://vitejs.dev/)
[![pnpm](https://img.shields.io/badge/pnpm-10-F69220?logo=pnpm&logoColor=white)](https://pnpm.io/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

## Tech Stack

| Category              | Technology              |
| --------------------- | ----------------------- |
| **Frontend (Public)** | Lit, Vite, TypeScript   |
| **Frontend (Admin)**  | React, Vite, TypeScript |
| **Backend (API)**     | Go, Air (Live reload)   |
| **Infrastructure**    | Docker, Docker Compose  |
| **Package Manager**   | pnpm (Workspace)        |

## Local Dev

### 0. セットアップ

```sh
mise i
```

### 1. コンテナのビルドと起動

```sh
docker compose build
docker compose up -d
```

### 2. バックエンド (Go API) の起動

Air を使用してホットリロード付きでAPIサーバーを起動。

```sh
air \
  --proxy.proxy_port "${API_PORT}" \
  -c ./backend/.air.toml \
  | jq .
```

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
