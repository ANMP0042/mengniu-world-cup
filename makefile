.PHONY: start
start:
	@${TOOLS_SHELL} go build && supervisorctl start mengniu:*
	@echo start
.PHONY: stop
stop:
	@${TOOLS_SHELL} supervisorctl stop mengniu:*
	@echo
.PHONY: run
run:
	@${TOOLS_SHELL} go run main.go
	@echo
