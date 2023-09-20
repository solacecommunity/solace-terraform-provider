## https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: help

default: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build: ### Build Solace Terraform Provider
	go build

