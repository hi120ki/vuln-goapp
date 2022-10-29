all: help

help: ## Print this help message
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: ping
ping: ## ping
	http http://localhost:8080/ping

.PHONY: ace
ace: ## Arbitrary Code Execution : whoami
	http POST http://localhost:8080/ace arg="whoami"

.PHONY: read
read: ## read file : /etc/passwd
	http POST http://localhost:8080/file/read arg="/etc/passwd"

.PHONY: create
create: ## create file : /app/tmp
	http POST http://localhost:8080/file/create arg="/app/tmp"

.PHONY: append
append: ## append file : /app/tmp
	http POST http://localhost:8080/file/create arg="/app/tmp"
	http POST http://localhost:8080/file/append arg="/app/tmp" arg2="test"
	http POST http://localhost:8080/file/read arg="/app/tmp"

.PHONY: delete
delete: ## delete file : /app/tmp
	http POST http://localhost:8080/file/create arg="/app/tmp"
	http POST http://localhost:8080/file/delete arg="/app/tmp"

.PHONY: download
download: ## download file : /app/tmp
	http POST http://localhost:8080/file/download arg="https://hi120ki.github.io/blog/robots.txt"
	cat robots.txt
	rm robots.txt

.PHONY: exec
exec: ## exec file : /app/tmp
	http POST http://localhost:8080/file/download arg="https://github.com/hi120ki/ctf4b-2022-testbin/raw/main/hello"
	http POST http://localhost:8080/file/chmod arg="hello"
	http POST http://localhost:8080/ace arg="./hello"

.PHONY: get
get: ## http get : http://httpbin.org/get
	http POST http://localhost:8080/http/get arg="http://httpbin.org/get"

.PHONY: json
json: ## http post : http://httpbin.org/post
	http -f POST http://localhost:8080/http/json arg="http://httpbin.org/post" arg2='{"token":"abcd"}'
