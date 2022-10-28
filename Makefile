all: help

help: ## Print this help message
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: ping
ping: ## ping
	http http://localhost:8080/ping

.PHONY: ace
ace: ## Arbitrary Code Execution : whoami
	http http://localhost:8080/ace arg="whoami"

.PHONY: read
read: ## Reading arbitrary files : /etc/passwd
	http http://localhost:8080/file/read arg="/etc/passwd"

.PHONY: create
create: ## Create file : /etc/passwd
	http http://localhost:8080/file/create arg="/etc/passwd"

.PHONY: append
append: ## append : /app/tmp
	http http://localhost:8080/file/create arg="/app/tmp"
	http http://localhost:8080/file/append arg="/app/tmp" arg2="test"
	http http://localhost:8080/file/read arg="/app/tmp"

.PHONY: delete
delete: ## delete : /app/tmp
	http http://localhost:8080/file/create arg="/app/tmp"
	http http://localhost:8080/file/delete arg="/app/tmp"
