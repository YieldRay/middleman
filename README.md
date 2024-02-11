# middleman

基于 HTTP 代理的，HTTP/HTTPS 拦截器

## Usage

CA 使用本地生成的自签名证书及密钥，注意必须系统信任自签名证书

```sh
middleman --help
```

要搭建代理服务器，参见 [此处示例](./server_examples)

## Build

命令行

```sh
go build
```

仅 GUI

```sh
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os windows
```
