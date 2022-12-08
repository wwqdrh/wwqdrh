SHELL := /bin/bash -o pipefail

.PHONY: install-easytest
install-easytest:
	wget -o ./vendor/bin/easytest https://github.com/wwqdrh/easytest/releases/download/v0.0.1/easytest-linux-amd64
	chmod +x ./vendor/bin/easytest

.PHONY: how_api_test
how_api_test:
	@./vendor/bin/easytest -json ./how_api_test/test.json
