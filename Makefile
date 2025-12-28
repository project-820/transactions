ifneq (,$(wildcard .env))
include .env
endif

GOOSE_DRIVER := postgres
GOOSE_DIR    := $(CURDIR)/internal/adapters/db/postgres/migrations

NAMESPACE := project820
SERVICE := transactions-service
BUILD_PATH := ./build
K8S_PATH := ./deployments/k8s

BUF=buf
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf


.PHONY: build apply restart logs proto clean vendor-googleapis vendor-protoc-gen-openapiv2 migrate-up migrate-down

build:
	eval $$(minikube docker-env); \
	docker build -t $(SERVICE):latest -f $(BUILD_PATH)/Dockerfile .

apply:
	kubectl apply -f $(K8S_PATH)/namespace.yaml
	kubectl apply -f $(K8S_PATH)/deployment_with_probes.yaml
	kubectl apply -f $(K8S_PATH)/service_cluster_ip.yaml
	kubectl apply -f $(K8S_PATH)/ingress.yaml;

restart:
	kubectl rollout restart deployment $(SERVICE) -n $(NAMESPACE) || true

logs:
	kubectl logs -n $(NAMESPACE) -l app=$(SERVICE)

vendor-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/googleapis &&\
	cd $(VENDOR_PROTO_PATH)/googleapis &&\
	git checkout
	mv $(VENDOR_PROTO_PATH)/googleapis/google $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/googleapis

vendor-protoc-gen-openapiv2:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_PATH)/grpc-gateway && \
 	cd $(VENDOR_PROTO_PATH)/grpc-gateway && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_PATH)/grpc-gateway/protoc-gen-openapiv2/options $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_PATH)/grpc-gateway

proto:
	$(BUF) generate --template ./buf.gen.yaml --path ./proto/api/transactions.proto

clean:
	rm -rf /pkg/*

migrate-up:
	goose -dir $(GOOSE_DIR) \
		$(GOOSE_DRIVER) \
		"user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DB) host=$(PG_HOST) port=$(PG_PORT) sslmode=$(PG_SSLMODE)" up

migrate-down:
	goose -dir $(GOOSE_DIR) \
		$(GOOSE_DRIVER) \
		"user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DB) host=$(PG_HOST) port=$(PG_PORT) sslmode=$(PG_SSLMODE)" down
