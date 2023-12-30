set -eu
cd ..
BUILD_PATH=.ci/build/compositor
rm -rf $BUILD_PATH
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=1
export CC_FOR_TARGET=gcc-aarch64-linux-gnu
export CC=aarch64-linux-gnu-gcc
go build -o=$BUILD_PATH
cp ./app/tools/scripts/* .ci/build/
cp prod-settings.yml .ci/build/settings.yml
cp ports-storage.db .ci/build/
