#!/bin/sh
rm -rf ./release
mkdir  release
go build -o chat
chmod +x ./chat
cp chat ./release/
#cp favicon.ico ./release/
cp -rf ./asset ./release/
cp -rf ./views ./release/