GO ?= go

run-all:
	@$(MAKE) -j3 run-command run-query run-worker

run-command:
	$(GO) -C packages/order-command run .

run-query:
	$(GO) -C packages/order-query run .

run-worker:
	$(GO) -C packages/order-worker run .
