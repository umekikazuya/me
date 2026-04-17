# me

## 開発環境操作

```sh
docker compose exec -T build
docker compose exec -T up -d
air --proxy.proxy_port ${API_PORT} -c ./backend/.air.toml | jq .
```
