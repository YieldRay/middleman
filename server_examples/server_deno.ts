//@ts-ignore
Deno.serve(async (request: Request) => {
    const url = new URL(request.url);

    if (url.pathname.startsWith("/http"))
        // proxy response
        return fetch(url.pathname.slice(1) + url.search, {
            method: request.method,
            headers: request.headers,
            body: request.body,
            redirect: "manual",
        }).catch((e: TypeError) => new Response(e.message, { status: 599, statusText: "Fail to fetch" }));

    // fake response
    return Response.json({
        method: request.method,
        url: request.url,
        headers: Object.fromEntries(request.headers.entries()),
        body: await request.text(),
    });
});
