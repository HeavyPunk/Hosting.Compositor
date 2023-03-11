BUILD_PATH=./build/compositor
rm -rf $BUILD_PATH
go build -o=$BUILD_PATH
cp ./app/tools/scripts/* ./build/
cp prod-settings.yml ./build/settings.yml
cp ports-storage.db ./build/