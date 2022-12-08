SHELL := /bin/bash -o pipefail

.PHONY: build
build:
	cd $(CURDIR)/appv1 && docker build -t grayapp:v1 .
	cd $(CURDIR)/appv2 && docker build -t grayapp:v2 .
	cd $(CURDIR)/gateway && docker build -t grayapp_gateway .

.PHONY: deploy
deploy:
	docker stack deploy -c stack.yaml gray

.PHONY: clean
clean:
	docker stack rm gray
