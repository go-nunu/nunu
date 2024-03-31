# 非容器部署（PM2+Nginx）
## PM2

pm2 是一个进程管理器，可以管理后端服务进程，还可以可以监控进程状态、查看进程日志以及资源占用情况等信息。
**安装PM2**
由于`PM2`依赖于`nodejs`，所以需要在你的服务器先安装`nodejs`。

然后直接使用`npm`安装`PM2`即可

```
npm i -g install pm2
```


在项目根目录，创建pm2配置文件`pm2.json`
```json
{
  "apps": [
    {
      "name": "nunu-api",
      "script": "./bin/server",
      "instances": 1,
      "exec_mode": "fork",
      "args":  "-conf=config/prod.yml"
    }
  ]
}
```
启动服务
```
pm2 start pm2.json
```

## Nginx配置
```nginx
server {
    listen 80;
	listen 443 ssl http2;
    server_name xxx.com;   # 填写你的域名
    index index.php index.html index.htm default.php default.htm default.html;
    root /data/www/wwwroot/xxx.com; # 替换为你的真实部署路径

    #error_page 404/404.html;
    ssl_certificate    /xxx/fullchain.pem;      # ssl证书路径
    ssl_certificate_key    /xxx/privkey.pem;    # ssl证书路径
    ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
    ssl_ciphers EECDH+CHACHA20:EECDH+CHACHA20-draft:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    location / {
      proxy_pass http://127.0.0.1:8000; # 修改为你自己的go服务端口
      proxy_set_header    Host             $host:$server_port;
      proxy_set_header    X-Real-IP        $remote_addr;
      proxy_set_header    X-Forwarded-For  $proxy_add_x_forwarded_for;
      proxy_set_header    HTTP_X_FORWARDED_FOR $remote_addr;
      proxy_redirect      off;
      proxy_buffering off;
    }
}
```