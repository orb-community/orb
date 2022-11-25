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
    while not event.is_set() and monitored_duration < int(monitor_time *  3600):
        container = client.containers.get(context.container_id)
        print("Monitor docker stats")
        statistics = container.stats(stream=False)
        used_memory = int(statistics["memory_stats"]["usage"]) - int(statistics["memory_stats"]["stats"]["file"])
        available_memory = int(statistics["memory_stats"]["limit"])
        cpu_delta = int(statistics["cpu_stats"]["cpu_usage"]["total_usage"]) - int(statistics["precpu_stats"]["cpu_usage"]["total_usage"])
        print("Memory Usage", (used_memory / available_memory) * 100.0, " %")
        print("CPU Delta ", cpu_delta)
        print("Memory Stats", statistics["memory_stats"])
        print("CPU Stats", statistics["cpu_stats"])
        print("PreCPU Stats", statistics["precpu_stats"])
        event.wait(30)
        monitored_duration = datetime.now().timestamp() - started
