test:
	poetry run locust --headless --users 50000 --spawn-rate 100 -H http://localhost:8080 -t 30s --only-summary
	# poetry run locust --headless --users 50000 --spawn-rate 50 -H http://localhost:8080
