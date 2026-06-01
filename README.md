# me

## 開発環境操作

```sh
docker compose build
docker compose up -d
air --proxy.proxy_port "${API_PORT}" -c ./backend/.air.toml | jq .
```

別ドメイン構成の動作確認をローカルで行う場合は、backend 起動時に `CORS_ALLOWED_ORIGINS` へ許可する frontend origin をカンマ区切りで渡します。これは CORS だけでなく、unsafe method に対する app 側の Origin 検証にも使われます。

```sh
CORS_ALLOWED_ORIGINS="http://localhost:5173" air --proxy.proxy_port "${API_PORT}" -c ./backend/.air.toml | jq .
```
