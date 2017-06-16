from locust import HttpLocust, TaskSet, task

class FailserverTaskSet(TaskSet):
    @task
    def index(self):
        self.client.get("/")

class FailserverLocust(HttpLocust):
    task_set = FailserverTaskSet
    min_wait = 100
    max_wait = 180
