# proxy server examples

To host a proxy server that can run with `middleman proxy <proxy-url>`,  
a HTTP server is required to be served at `proxy-url`, for example `https://my-server.com`,  
the server should forward the request method, headers and body to target url specified by the requested path (without the prefix `/`)

For example, the request may look like this:

```
HTTP/1.1 POST /https://example.net
Host: my-server.com
Accept-Encoding: gzip, deflate, br

some_request_body
```

This directory contains some example implement in various programming languages.
