.PHONY: help init_db local_deploy, db, dash, clear-mlflow mlflow, export_req, export_req_dev clear test heroku-deploy heroku-log

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

db: ## Sobe banco de dados
	docker-compose up --build db

test: ## Realiza testes
	cd fetch && go test
