# The Donders Repository API gateway

This API gateway provides RESTful API interfaces for querying and searching on a limited set of Donders Repository collection and user attributes.

## how it works

Every 10 minutes, it queries all collections and users available to the RDR service account (defined as `irodsUser` and `irodsPass` attributes in [config.yml](config/config.yml) ).  The results are cached in memory for quick response to the API client. When interacting with iRODS iCAT service, it makes use of the [go-irodsclient](https://github.com/cyverse/go-irodsclient).

The API interface is defined by [swagger.yml](pkg/swagger/swagger.yaml).

## available collection attributes

The following collection attributes are queried from iCAT:

- Identifier
- IdentifierDOI
- ProjectID
- Type
- State
- OrganisationalUnit
- QuotaInBytes
- SizeInBytes
- NumberOfFiles

See [here](https://github.com/dccn-tg/dr-gateway/blob/e9fb2cd0c63b2c0fa72bcdd0fdd2d8da212d2cfd/pkg/dr/collection.go#L94).

## available user attributes

The following user attributes are queried from iCAT:

- DisplayName
- IdentityProvider
- Email
- OrganizationalUnit

See [here](https://github.com/dccn-tg/dr-gateway/blob/e9fb2cd0c63b2c0fa72bcdd0fdd2d8da212d2cfd/pkg/dr/user.go#L13).

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
