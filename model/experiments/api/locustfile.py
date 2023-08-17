from locust import HttpUser, task


class Locust(HttpUser):
    @task
    def pred_5000(self):
        self.client.post(
            "/predict", json={'input': [[220, 3, 4, 4, 3, 3] for _ in range(5_000)]})
