# Findings of HTTP protocol and Host header

TODO: finish this!

**TL;DR:** in HTTP/1.1, a Host header is required, which is how the same IP address can service multiple DNS names while being able to distinguish between them.

I ran a handful of tests using `netcat` and httpbin.org. httpbin.org uses nodejs and expressjs.

`Netcat` is a utility that allows you to easily handcraft HTTP requests by connecting directly to a web server TCP socket and simply typing in the request.
The connection is then terminated with an EOF character (Ctrl-z on Windows, Ctrl-d on Mac/Linux).

## HTTP protocol and Host header

When a hostname is not specified, the implicit URL appears to be the server's DNS name (or presumably its IP address if it doesn't have a DNS name).
I imagine this might be dependent on the server.

The hostname can also be included in the first line of the HTTP request. For 0.9 and 1.0, this works in lieu of a Host header.
In 1.1, a Host header seems to be required.

### Note on HTTP 1.1 _with_ a full URL in the request but _without_ a Host header

This results in _400 Bad Request_.

### Summary of Host header tests

| Protocol version | Absolute URL in request? | Header required? | Response Headers         | Request URL              |
| ---------------- | ------------------------ | ---------------- | ------------------------ | ------------------------ |
| 0.9              | No                       | No               | None and/or not required | Inferred from DNS name   |
| 0.9              | Yes                      | No               | None and/or not required | Specified in request     |
| 1.0              | No                       | No               | Present and/or required  | Inferred from DNS name   |
| 1.0              | Yes                      | No               | Present and/or required  | Specified in request     |
| 1.1              | No                       | Required         | Present and/or required  | Specified in Host header |

## HTTP 0.9

### Request

No HTTP protocal specified, therefore HTTP 0.9.

- Host header is not required.

### Results

- Server response has no headers. They're either not required or not supported.
- _The request URL claimed by the web server is the Amazone hostname, combined with `/get` specified in the HTTP request._

```bash
$ nc httpbin.org 80
GET /get
{
  "args": {},
  "headers": {
    "Host": "a0207c42-pmhttpbin-pmhttpb-c018-592832243.us-east-1.elb.amazonaws.com",
    "X-Amzn-Trace-Id": "Root=1-5f921406-6725732b5e7727ba21c487e5"
  },
  "origin": "45.23.234.8",
  "url": "http://a0207c42-pmhttpbin-pmhttpb-c018-592832243.us-east-1.elb.amazonaws.com/get"
}
```

## HTTP 1.0

### Request

- Host header is not required.

### Response

- Perhaps response headers are required by to be sent by the server.
- _The request URL claimed by the web server is the Amazone hostname, combined with `/get` specified in the HTTP request._

```bash
$ nc httpbin.org 80
GET /get HTTP/1.0

HTTP/1.1 200 OK
Date: Thu, 22 Oct 2020 23:21:53 GMT
Content-Type: application/json
Content-Length: 312
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "headers": {
    "Host": "a0207c42-pmhttpbin-pmhttpb-c018-592832243.us-east-1.elb.amazonaws.com",
    "X-Amzn-Trace-Id": "Root=1-5f921411-54a5eec444ca2761780c4059"
  },
  "origin": "45.23.234.8",
  "url": "http://a0207c42-pmhttpbin-pmhttpb-c018-592832243.us-east-1.elb.amazonaws.com/get"
}
```

## Remainder -- to be completed

```bash
# HTTP 1.1, which is historically the dominant version.
# Host header apparently required
➜  ~ nc httpbin.org 80
GET /get HTTP/1.1

HTTP/1.1 400 Bad Request
Server: awselb/2.0
Date: Thu, 22 Oct 2020 23:22:30 GMT
Content-Type: text/html
Content-Length: 122
Connection: close

<html>
<head><title>400 Bad Request</title></head>
<body>
<center><h1>400 Bad Request</h1></center>
</body>
</html>

# HTTP 1.1 with a Host header
# Note that the URL claimed by the server is a combination of the Host header and the URL specified in the GET
$ nc httpbin.org 80
GET /get HTTP/1.1
Host: httpbin.org

HTTP/1.1 200 OK
Date: Thu, 22 Oct 2020 23:22:47 GMT
Content-Type: application/json
Content-Length: 196
Connection: keep-alive
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "headers": {
    "Host": "httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-5f921447-1cdba2c70e0b2a116eab7eaf"
  },
  "origin": "45.23.234.8",
  "url": "http://httpbin.org/get"
}
```

## Results round 2, with full URL in request

TO DO: finish this

```bash
➜  ~ nc httpbin.org 80
GET http://httpbin.org/get
{
  "args": {},
  "headers": {
    "Host": "httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-5f921c42-65737d0f15b51e8a2a8dd5cf"
  },
  "origin": "45.23.234.8",
  "url": "http://httpbin.org/get"
}
GET http://httpbin.org/get HTTP/1.0

Ncat: Broken pipe.
➜  ~ nc httpbin.org 80
GET http://httpbin.org/get HTTP/1.0

HTTP/1.1 200 OK
Date: Thu, 22 Oct 2020 23:57:16 GMT
Content-Type: application/json
Content-Length: 196
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "headers": {
    "Host": "httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-5f921c5c-725b81754c277214501eabef"
  },
  "origin": "45.23.234.8",
  "url": "http://httpbin.org/get"
}
➜  ~ nc httpbin.org 80
GET http://httpbin.org/get HTTP/1.0
^C
➜  ~ nc httpbin.org 80
GET http://httpbin.org/get HTTP/1.1

HTTP/1.1 400 Bad Request
Server: awselb/2.0
Date: Thu, 22 Oct 2020 23:57:40 GMT
Content-Type: text/html
Content-Length: 122
Connection: close

<html>
<head><title>400 Bad Request</title></head>
<body>
<center><h1>400 Bad Request</h1></center>
</body>
</html>
➜  ~


```
