
setup-env:
	python3 -m venv behave_orb
	source behave_orb/bin/activate
	pip install -r requirements.txt

behave:
	behave

behave-json:
	behave -f json.pretty

behave-xml:
	behave --junit

install-allure:
	pip install allure-behave
	allure serve allure_server_orb

behave-allure:
	behave -f allure_behave.formatter:AllureFormatter –o allure_server_orb
