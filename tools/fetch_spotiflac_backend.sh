#!/bin/sh

set -e

FOLDER=lib

rm -rf lib/
git clone https://github.com/afkarxyz/SpotiFLAC.git
cp -r SpotiFLAC/backend/ lib/
mkdir -p app/ 
cp SpotiFLAC/app.go app/app.go
rm -rf SpotiFLAC

sed -i "s/package main/package app/g" app/app.go
sed -i "s/\"spotiflac\/backend\"/backend \"github.com\/Superredstone\/spotiflac-cli\/lib\"/g" app/app.go

# Nix shenanigans
chmod -R 777 lib

for i in $(ls lib/); do
	sed -i "s/package backend/package $FOLDER/g" $FOLDER/$i
done

