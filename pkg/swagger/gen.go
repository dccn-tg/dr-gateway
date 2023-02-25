package swagger

//go:generate rm -rf server/models server/restapi
//go:generate mkdir -p server
//go:generate swagger generate server --quiet --target server --name dr-gateway --spec swagger.yaml --exclude-main -P models.Principal
//go:generate rm -rf client/models client/client
//go:generate mkdir -p client
//go:generate swagger generate client --quiet --target client --name dr-gateway --default-scheme=https --spec swagger.yaml
