# emoine
これはemoineもどき

## 使い方
`/etc/hosts`に`127.0.0.1 emoine.trapti.tech`を追加
```
task up
cd cmd/client
go run main.go
```
streamに接続している間はコメントを送信できないから、別のターミナルから送る必要がある  
`localhost:8080`に接続すると、traefikのダッシュボード画面を見れる

## サーバーをスケールアウトする
dockerなら
```
docker compose up -d --scale server=2
```
これで増えます