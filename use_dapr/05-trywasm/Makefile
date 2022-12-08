SHELL := /bin/bash -o pipefail

.PHONY: basic
basic:
	cd $(CURDIR)/basic && cargo build --release && ./target/release/basic_hyper_server

.PHONY: wasi
wasi:
	cd $(CURDIR)/wasi && cargo build --target wasm32-wasi --release && ~/.wasmedge/bin/wasmedge target/wasm32-wasi/release/wasmedge_hyper_server.wasm