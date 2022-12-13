```markdown
docker stack deploy -c stack.yaml test

docker service create -p 8000:8000 --mount type=bind,src=`pwd`/stream.conf,dst=/etc/nginx/conf.d/basic.conf --mount type=bind,src=`pwd`/lib,dst=/usr/local/openresty/nginx/conf/lua_modules/  openresty/openresty

curl http://192.168.0.111:8000/api1

docker service logs test_auth

docker service logs test_basic

docker service logs test_basic_ingress

docker service update --force test_basic_ingress
```
