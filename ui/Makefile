ui:
	docker build --tag=mainflux/ui -f docker/Dockerfile .

ui-experimental:
	docker build --tag=mainflux/ui-experimental -f docker/Dockerfile.experimental .

run:
	docker-compose -f docker/docker-compose.yml up

clean:
	docker-compose -f docker/docker-compose.yml down --rmi all -v --remove-orphans

release:
	$(eval version = $(shell git describe --abbrev=0 --tags))
	git checkout $(version)
	docker tag mainflux/ui mainflux/ui:$(version)
	docker push mainflux/ui:$(version)
