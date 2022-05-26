[test_config]
# Required fields:
email=<user-email>
password=<user-password>
orb_address=<orb url>
prometheus_username= <Your Grafana Cloud Prometheus username>
prometheus_key= <Your Grafana Cloud API Key. Be sure to grant the key a role with metrics push privileges>
remote_prometheus_endpoint= <base URL to send Prometheus metrics to Grafana Cloud>

This field is required if use docker approach to run the tests
# orb_path = <path to orb directory>

# Optional fields:
# agent_docker_image=ns1labs/orb-agent
# agent_docker_tag=latest
# orb_agent_interface=mock
# ignore_ssl_and_certificate_errors=True
# is_credentials_registered=False
