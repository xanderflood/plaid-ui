#!/usr/bin/env bash
set +e

docker login -u $DOCKER_USER -p $DOCKER_PASS

# build the executable
#
#
# TODO build with go 1.13 in a docker container?
#.  or is there travis config for that?
#
#
#
docker run --rm -v $(pwd):/builddir -w /builddir -e CGO_ENABLED=0 -e GOOS=linux golang:1.13 go build -v -o ./build/api/api ./cmd/api
# go version
# CGO_ENABLED=0 GOOS=linux go build -o ./build/api/api ./cmd/api

export docker_repo=xanderflood/plaid-ui
export docker_build_directory=./build/api

# build and push the base image
export TAG=`if [ "$TRAVIS_BRANCH" = "master" ]; then echo "latest"; else echo "staging" ; fi`
export tmpname="local"
export dockerfile="$docker_build_directory/Dockerfile"
export tags="build-${TRAVIS_BUILD_NUMBER},commit-${TRAVIS_COMMIT::8},$TAG"
./build_and_push_image.sh

# build and push the swarm-focused image
export TAG=`if [ "$TRAVIS_BRANCH" = "master" ]; then echo "swarm"; else echo "swarm-staging" ; fi`
export tmpname="swarm-local"
export dockerfile="$docker_build_directory/swarm.Dockerfile"
export tags="build-${TRAVIS_BUILD_NUMBER},commit-${TRAVIS_COMMIT::8},$TAG"
./build_and_push_image.sh
