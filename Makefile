SHELL=/bin/bash
.DEFAULT_GOAL := help

.PHONY: print-order-book
print-order-book: ## Print order book.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/orderbook/main.go

.PHONY: print-orders
print-orders: ## Print orders.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/order/fetch/main.go

.PHONY: create-orders
create-order: ## Create order.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/order/create/main.go

.PHONY: count-go
count-go: ## Count number of lines of all go codes.
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

# See "Self-Documented Makefile" article
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
