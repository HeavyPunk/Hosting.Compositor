set -eu

bash build.sh

DOCKER_IMAGE=kirieshki/simple-hosting.compositor
DOCKER_TAG=2024.01.17-amd64

docker buildx build -t $DOCKER_IMAGE:$DOCKER_TAG --progress plain --platform linux/amd64 --push .
