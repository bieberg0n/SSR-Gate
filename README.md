# SSR-Gate
根据ssr订阅链接自动连接可用的服务器

---

## Usage

```shell script
wget https://github.com/bieberg0n/SSR-Gate/releases/download/0.0.1/SSR-Gate

./SSR-Gate -h
Usage of ./SSR-Gate:
  -h	help
  -l int
    	listen port (default 1080)
  -u string
    	ssr url

./SSR-Gate -u "https://your-ssr-link/balabala"
```

Then new term window:
```
curl --proxy socks5://127.0.0.1:1080 http://google.com                                                                                                                                       (15:11:55)───┘
<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>301 Moved</TITLE></HEAD><BODY>
<H1>301 Moved</H1>
The document has moved
<A HREF="http://www.google.com/">here</A>.
</BODY></HTML>
```
