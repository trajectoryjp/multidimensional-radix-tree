# クロスビルド

## ターゲット

GOがセットアップされている環境で下記のコマンドを実行するとサポートされているターゲットが表示される

```
go tool dist list
```

## 具体例

- Raspberry pi

```
GOOS=linux GOARCH=arm go build -o trjx_transfer_raspberry
```

- Raspberry pi (1,zero)
```
GOOS=linux GOARCH=arm GOARM=6 go build -o trjx_transfer_raspberry_zero
```
https://qiita.com/m0a/items/d933982293dcadd4998c

- Raspberry pi (2)
```
GOOS=linux GOARCH=arm GOARM=7 go build -o trjx_transfer_raspberry_2
```
https://qiita.com/m0a/items/d933982293dcadd4998c

- centos

```
GOOS=linux GOARCH=amd64 go build -o trjx_transferlinux
```

- intel Macintosh

```
GOOS=darwin GOARCH=amd64 go build -o trjx_transfer_intel_mac
```

- Apple Silicon Macintosh

```
GOOS=darwin GOARCH=arm64 go build -o trjx_transfer_m_mac
```

- Ubuntu AMD

```
GOOS=linux GOARCH=amd64 go build -o trjx_transfer_linux_amd
```


# 設定

実行ファイルと同じ階層の`setting`ディレクトリ内の`config.json`により動作を切り替えることができる。
