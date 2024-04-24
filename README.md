# go实现高并发秒杀商城

## 1. 环境搭建
- rabbitmq
    - `docker run -it -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3 /bin/bash` 下载rabbitmq的docker, 在container中跑rabbitmq
    - 在container中: `rabbitmq-server &`
    - 浏览器中访问http://127.0.0.1:15672/#/, 即可进入管理界面. 默认账号密码guest/guest