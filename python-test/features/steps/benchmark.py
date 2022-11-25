import time

import docker
from behave import then, step
client = docker.from_env()
from utils import threading_wait_until


@step("monitor the activity of memory usage during {monitor_time} minutes")
def monitor_docker_stats_during(context, monitor_time):
    monitored_duration = 0
    while monitored_duration >= int(monitor_time):
        container = client.containers.get(context.container_id)
        print_stats_and_wait(container, wait_time=30)
        monitored_duration += 0.5


@threading_wait_until
def print_stats_and_wait(container):
    print("Monitor docker stats")
    statistics = container.stats(stream=False)
    print(statistics)
    time.sleep(30)
