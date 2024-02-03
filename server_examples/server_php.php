<?php
error_reporting(E_ALL);
ini_set("display_errors", 1);
main();

function main()
{
    $path = urlPathPart();
    reverseProxy($path);
}

/**
 * =================
 * COMMON FUNCTIONS
 * =================
 */

// string utils

function startsWith($haystack, $needle)
{
    return substr_compare($haystack, $needle, 0, strlen($needle)) === 0;
}
function endsWith($haystack, $needle)
{
    return substr_compare($haystack, $needle, -strlen($needle)) === 0;
}
function removePrefix($haystack, $needle)
{
    if (!startsWith($haystack, $needle)) {
        return $haystack;
    }
    return substr($haystack, strlen($needle));
}
function removeSuffix($haystack, $needle)
{
    if (!endsWith($haystack, $needle)) {
        return $haystack;
    }
    return substr($haystack, 0, -strlen($needle));
}


/**
 * get fixed url path part (no prefix '/')
 */
function urlPathPart()
{
    $PATH = "";
    $REQUEST_URI = $_SERVER["REQUEST_URI"];
    $SCRIPT_NAME = $_SERVER["SCRIPT_NAME"];

    if (startsWith($REQUEST_URI, $SCRIPT_NAME)) {
        // no rewrite
        // $_SERVER['REQUEST_URI']	/test/index.php/https://example.net
        // $_SERVER['SCRIPT_NAME']	/test/index.php
        $PATH = removePrefix(removePrefix($REQUEST_URI, $SCRIPT_NAME), "/");
    } else {
        // has rewrite
        // $_SERVER['REQUEST_URI']	/test/https://example.net
        // $_SERVER['SCRIPT_NAME']	/test/index.php
        $PATH = removePrefix($REQUEST_URI, removeSuffix($SCRIPT_NAME, "index.php"));
    }
    /**
     * <IfModule mod_rewrite.c>
     * RewriteEngine On
     * RewriteBase /test/
     * RewriteCond %{REQUEST_FILENAME} !-f
     * RewriteCond %{REQUEST_FILENAME} !-d
     * RewriteRule ^(.*)$ index.php [L,E=PATH_INFO:$1]
     * </IfModule>
     */
    return $PATH;
}

/**
 * reverse proxy for an url
 * 
 * @param targetUrl  -  the url that proxy to
 */
function reverseProxy($targetUrl)
{
    // Get incoming request headers
    foreach (getallheaders() as $key => $val) {
        // Exclude some header
        if (strtolower($key) !== "host" && strtolower($key) !== "accept-encoding") {
            $requestHeaders[] = "$key: $val";
        }
    }

    // Initialize cURL session
    $ch = curl_init();

    // Set the target URL
    curl_setopt($ch, CURLOPT_URL, $targetUrl);

    // Enable automatically set the Referer: field
    curl_setopt($ch, CURLOPT_AUTOREFERER, true);

    // Pass incoming request headers to the target server
    curl_setopt($ch, CURLOPT_HTTPHEADER, $requestHeaders);

    // Forward the request method and body
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $_SERVER["REQUEST_METHOD"]);
    curl_setopt($ch, CURLOPT_POSTFIELDS, file_get_contents("php://input"));

    // Do not follow location
    curl_setopt($ch, CURLOPT_FOLLOWLOCATION, false);

    // Do not fail on error
    curl_setopt($ch, CURLOPT_FAILONERROR, false);

    // Forward headers
    curl_setopt($ch, CURLOPT_HEADERFUNCTION, function ($curl, $header_line) {
        if (startsWith($header_line, "http/") || startsWith($header_line, "transfer-encoding")) {
            // skip http version header
            return strlen($header_line);
        }
        http_response_code(curl_getinfo($curl, CURLINFO_HTTP_CODE)); // forward status code
        header($header_line); // forward single header
        return strlen($header_line); // curl need this return
    });

    // Skip ssl check
    curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, false);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
    curl_setopt($ch, CURLOPT_SSL_VERIFYSTATUS, false);
    curl_setopt($ch, CURLOPT_PROXY_SSL_VERIFYPEER, false);

    // Execute the cURL request & send response
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, false);
    $success = curl_exec($ch);

    // Error handler
    if (!$success) {
        http_response_code(599);
        header("content-type: application/json");
        $errno = curl_errno($ch);
        $error = curl_error($ch);
        $json =  json_encode(["errno" => $errno, "error" => $error]);
        if ($_SERVER["REQUEST_METHOD"] === "HEAD") header("x-proxy-error: $json"); // experimental
        else echo $json;
    }
    exit();
}
