# dr-gateway

API gateway of the Donders Repository.

## build the docker container

```bash
$ docker-compose build --force-rm
```

## run the docker container

```bash
$ cp env.sh .env
$ docker-compose up -d
```

or run the [start script](start.sh)

```bash
$ ./start.sh
```

## API document

Once the service is up and running, the API document can be found at the endpoint `/v1/docs`.
