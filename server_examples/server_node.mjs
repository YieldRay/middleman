import http from "http";
import https from "https";
import url from "url";

const PORT = Number(process.env.PORT) || 3000;

const server = http.createServer((req, res) => {
    const u = req.url.slice(1);
    const { protocol, hostname, port, path } = url.parse(u);
    const isHTTPS = protocol === "https:";
    const client = isHTTPS ? https : http;

    delete req.headers.host;

    const request = client.request(
        {
            hostname,
            port: Number(port) || (isHTTPS ? 443 : 80),
            path,
            method: req.method,
            headers: req.headers,
        },
        (response) => {
            res.writeHead(response.statusCode, response.statusMessage, response.headers);
            response.pipe(res);
        }
    );

    req.pipe(request);

    request.on("error", (e) => req.destroy(e));
});

server.listen(PORT).on("listening", () => console.log(`Server listen at http://localhost:${PORT}`));
