from datetime import datetime
import threading
import docker
from behave import then, step
import json

MINUTE_INT = 3600
client = docker.from_env()


@step("agent will be monitored in cpu and memory usage with {policy_amount} policies applied")
def prepare_monitor(context,policy_amount):
    context.policy_amount = policy_amount


@step("monitor the activity of memory usage during {monitor_time} minutes")
def monitor_docker_stats_during(context, monitor_time):
    started = datetime.now().timestamp()
    monitored_duration = 0
    event = threading.Event()
    sr_file_name = "benchamark/short-report-" + str(started) + ".json"
    short_report_file = open(sr_file_name, "w")
    lr_file_name = "benchamark/docker-container-stats-" + str(started) + ".json"
    long_report_file = open(lr_file_name, "w")
    while not event.is_set() and monitored_duration < int(monitor_time * MINUTE_INT):
        container = client.containers.get(context.container_id)
        print("Monitor docker stats")
        statistics = container.stats(stream=False)
        used_memory = int(statistics["memory_stats"]["usage"]) - int(statistics["memory_stats"]["stats"]["file"])
        available_memory = int(statistics["memory_stats"]["limit"])
        cpu_delta = int(statistics["cpu_stats"]["cpu_usage"]["total_usage"]) - int(statistics["precpu_stats"]["cpu_usage"]["total_usage"])
        # print to a file instead of printing in console
        print("Memory Usage", (used_memory / available_memory) * 100.0, " %")
        print("CPU Delta ", cpu_delta)
        short_report = {"timestamp": datetime.now().timestamp(), "cpu_delta": cpu_delta, "memory_usage": (used_memory / available_memory) * 100.0}
        json.dump(short_report, short_report_file)
        json.dump(statistics, long_report_file)
        print("Memory Stats", statistics["memory_stats"])
        print("CPU Stats", statistics["cpu_stats"])
        print("PreCPU Stats", statistics["precpu_stats"])
        event.wait(30)
        monitored_duration = datetime.now().timestamp() - started
