#!/bin/sh

APP="poulpe"
BUILD_FLAGS="-o ${APP}"
MODE=${MODE:-dev}
VERSION=development
test -f .version && VERSION=$(cat .version)

[ "${MODE}" = "prod" ] && BUILD_FLAGS="${BUILD_FLAGS}"
[ "${MODE}" = "dev" ] && BUILD_FLAGS="${BUILD_FLAGS} -v"

echo "Building ${APP}..."
CMD="go build ${BUILD_FLAGS}"
echo "Running ${CMD}"

eval ${CMD}
