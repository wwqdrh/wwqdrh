SHELL := /bin/bash -o pipefail

# 创建实验基础环境
.PHONY: init
init:
	docker network create --driver overlay test_logcenter

	docker volume create --driver local --label test test_logcenter_app1
	docker volume create --driver local --label test test_logcenter_app2
	docker volume create --driver local --label test test_logcenter_promtail_store
	docker volume create --driver local --label test test_logcenter_loki_store
	docker volume create --driver local --label test test_logcenter_grafana_store
	
	docker config create test_logcenter_loki $(CURDIR)/conf/loki.yaml
	docker config create test_logcenter_promtail $(CURDIR)/conf/promtail.yaml

.PHONY: clean
clean:
	docker network rm test_logcenter

	docker volume rm test_logcenter_app1
	docker volume rm test_logcenter_app2
	docker volume rm test_logcenter_promtail_store
	docker volume rm test_logcenter_loki_store
	docker volume rm test_logcenter_grafana_store

	docker config rm test_logcenter_loki
	docker config rm test_logcenter_promtail

.PHONY: test-logcenter
test-logcenter:
	docker stack deploy -c $(CURDIR)/stack.yaml logcenter

.PHONY: test-logcenter-clean
test-logcenter-clean:
	docker stack rm logcenter

# 程序A每隔5秒往挂载了volume的日志文件中写日志
# 程序B从挂载了同样volume的路径下去读日志，如果能够读到日志就说明成功了
.PHONY: test-volume-feature
test-volume-feature:
	docker service create \
	--name test_logcenter_volume_write \
	--mount type=volume,src=test_logcenter,dst=/data/log \
	--mount type=bind,src=$(CURDIR)/volume-write.py,dst=/app/volume-write.py python:3.10 python /app/volume-write.py

	docker service create \
	--name test_logcenter_volume_read \
	--mount type=volume,src=test_logcenter,dst=/data/log \
	--mount type=bind,src=$(CURDIR)/volume-read.py,dst=/app/volume-read.py python:3.10 python /app/volume-read.py

.PHONY: test-volume-feature-clean
test-volume-feature-clean:
	docker service rm test_logcenter_volume_read
	docker service rm test_logcenter_volume_write
