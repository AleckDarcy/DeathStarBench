#!/bin/bash

cd $(dirname $0)/..

CBPATH="/Users/aleck/go/src/github.com/AleckDarcy/ContextBus"

EXEC="docker buildx"

USER="aleck123"

TAG="latest"

# ENTER THE ROOT FOLDER
cd ../
ROOT_FOLDER=$(pwd)
$EXEC create --name mybuilder --use

rm vendor/github.com/AleckDarcy/ContextBus
cp -r $CBPATH vendor/github.com/AleckDarcy/ContextBus

for i in hotelreservation #frontend geo profile rate recommendation reserve search user #uncomment to build multiple images
do
  IMAGE=${i}
  echo Processing image ${IMAGE}
  cd $ROOT_FOLDER
  $EXEC build -t "$USER"/cb_dsb_"$IMAGE":"$TAG" -f Dockerfile . --platform linux/arm64,linux/amd64 --push
  cd $ROOT_FOLDER

  echo
done

rm -rf vendor/github.com/AleckDarcy/ContextBus
mkdir -p vendor/github.com/AleckDarcy/
ln -s $CBPATH vendor/github.com/AleckDarcy/ContextBus

cd - >/dev/null
