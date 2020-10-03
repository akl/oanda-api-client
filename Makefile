SHELL=/bin/bash
.DEFAULT_GOAL := help

.PHONY: print-order-book
print-order-book: ## Print order book.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/orderbook/fetch/main.go

.PHONY: print-order-book-vop
print-order-book-vop: ## Print order book vop.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/orderbook/vop/main.go

.PHONY: print-orders
print-orders: ## Print orders.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/order/fetch/main.go

.PHONY: create-order
create-order: ## Create order.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/order/create/main.go

.PHONY: cancel-order
cancel-order: ## Cancel order.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/order/cancel/main.go

.PHONY: close-trade
close-trade: ## Close trade.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/trade/close/main.go

.PHONY: print-trades
print-trades: ## Print trades.
	AWS_PROFILE=yukiinoue-private AWS_DEFAULT_REGION=ap-northeast-1 ENVIRONMENT=Practice go run cmd/trade/fetch/main.go


.PHONY: count-go
count-go: ## Count number of lines of all go codes.
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

# See "Self-Documented Makefile" article
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
