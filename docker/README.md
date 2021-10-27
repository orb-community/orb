# How to run locally


1 - On the root of this project run
```
make dockers_dev
```
to create the images

2 - Edit [docker/.env](docker/.env) if necessary

3 - Run docker-compose

```
docker-compose --env-file docker/.env -f docker/docker-compose.yml up -d
```

4 - Point the browser to http://localhost:80 (It's possible to change the port in the .env file)


