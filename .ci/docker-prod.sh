set -eu

bash build.sh

DOCKER_IMAGE=kirieshki/simple-hosting.compositor
DOCKER_TAG=latest

docker buildx build -t $DOCKER_IMAGE:$DOCKER_TAG --progress plain --platform linux/arm64/v8 --push .
