# middleman

基于 HTTP 代理的，HTTP/HTTPS 拦截器

## Usage

CA 使用本地生成的自签名证书及密钥，注意必须系统信任自签名证书

```sh
Usage:
  middleman [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inspect     Inspect http(s) traffic
  proxy       使用代理服务器，例如：https://cros.deno.dev/

Flags:
      --ca-crt string     CA证书路径，参考生成命令：openssl req -x509 -new -key ca.key -out ca.crt -days 3650 (default "ca.crt")
      --ca-key string     CA密钥路径，参考生成命令：openssl genpkey -algorithm RSA -out ca.key (default "ca.key")
      --expose            expose local server
  -h, --help              help for middleman
      --log               write log to file
      --log-level uint8   Set the log level, TRACE|DEBUG|INFO|WARN|ERROR|FATAL (default 2)
      --log-path string   path to log file of request (default "middleman_2024-02-02.log")
      --port int          http proxy local port (default 9980)

Use "middleman [command] --help" for more information about a command.
```

要搭建代理服务器，参见 [此处示例](./server_examples)
