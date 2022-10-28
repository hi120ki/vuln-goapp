all: help

help: ## Print this help message
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: ping
ping: ## ping
	http http://localhost:8080/ping

.PHONY: ace
ace: ## Arbitrary Code Execution : whoami
	http http://localhost:8080/ace arg="whoami"

.PHONY: raf
raf: ## Reading arbitrary files : /etc/passwd
	http http://localhost:8080/raf arg="/etc/passwd"
