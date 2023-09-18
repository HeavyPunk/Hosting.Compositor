set -eu
cd ..
BUILD_PATH=.ci/build/compositor
rm -rf $BUILD_PATH
go build -o=$BUILD_PATH
cp ./app/tools/scripts/* .ci/build/
cp prod-settings.yml .ci/build/settings.yml
cp ports-storage.db .ci/build/
