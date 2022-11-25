from datetime import datetime
import threading
import docker
from behave import then, step
client = docker.from_env()



@step("monitor the activity of memory usage during {monitor_time} minutes")
def monitor_docker_stats_during(context, monitor_time):
    started = datetime.now().timestamp()
    monitored_duration = 0
    event = threading.Event()
    while not event.is_set() and monitored_duration < int(monitor_time):
        container = client.containers.get(context.container_id)
        print("Monitor docker stats")
        statistics = container.stats(stream=False)
        print(statistics)
        event.wait(30)
        monitored_duration = datetime.now().timestamp() - started
