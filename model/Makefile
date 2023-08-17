local_deploy:
	poetry run uvicorn --host 0.0.0.0 --port 5060 "api:app" --reload

mlflow: ## Sobe MLFlow na porta 5050
	poetry run mlflow server --default-artifact-root mlruns -p 6050 --gunicorn-opts "--timeout 0"
