NAME    ?= go-angular-security
PORT    ?= 5000
ENVFILE ?= server/.env

.DEFAULT_GOAL := help

.PHONY: help build run logs remove replace

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the Docker image
	docker build -t $(NAME) .

run: ## Run the container (stops existing one first)
	-docker stop $(NAME) 2>/dev/null || true
	-docker rm $(NAME) 2>/dev/null || true
	docker run -d \
		--name $(NAME) \
		--env-file $(ENVFILE) \
		-p $(PORT):5000 \
		--restart unless-stopped \
		$(NAME)

logs: ## Tail container logs
	docker logs -f $(NAME)

remove: ## Stop and remove the container
	docker stop $(NAME)
	docker rm $(NAME)

replace: build run ## Rebuild image and restart container
