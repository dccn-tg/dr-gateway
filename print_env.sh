#!/bin/bash

echo "# version"
echo "DOCKER_IMAGE_TAG=$DOCKER_IMAGE_TAG"
echo 
echo "# docker registry endpoint"
echo "DOCKER_REGISTRY=$DOCKER_REGISTRY"
echo
echo "# configuration file for api-server"
echo "CFG_API_SERVER=$CFG_API_SERVER"
echo
echo '# service port for external client'
echo "DR_GATEWAY_EXTERNAL_PORT=$DR_GATEWAY_EXTERNAL_PORT"