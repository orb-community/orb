[test_config]
# Required fields:
email=<user-email>
password=<user-password>
orb_address=<orb url>

# Required fields for sink scenario:
# prometheus_username= <Your Grafana Cloud Prometheus username>
# prometheus_key= <Your Grafana Cloud API Key. Be sure to grant the key a role with metrics push privileges>
# remote_prometheus_endpoint= <base URL to send Prometheus metrics to Grafana Cloud> (ex. prometheus-prod-10-prod-us-central-0.grafana.net)

# Optional fields:
# agent_docker_image=ns1labs/orb-agent
# agent_docker_tag=latest
# orb_agent_interface=mock