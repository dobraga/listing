.PHONY: help db test config_heroku build deploy build_deploy front loca_deploy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

local_deploy: ## Deploy locals
	docker-compose up --build

db: ## Sobe banco de dados
	docker-compose up --build db

test: ## Realiza testes
	cd fetch && go test


config_heroku: ## Configura heroku
	heroku login
	heroku container:login

build: ## Builda imagem
	docker-compose build app
	docker tag listings_app registry.heroku.com/stark-castle-15501/web

deploy: ## Deploy heroku
	docker push registry.heroku.com/stark-castle-15501/web
	heroku container:release web

build_deploy: ## Build e deploy
	$(MAKE) build
	$(MAKE) deploy

front: ## Sobe front
	cd front && poetry run python main.py
