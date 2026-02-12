#!/bin/sh

set -e

FOLDER=lib
# USELESS_FILES="file_dialog.go folder.go"

rm -rf lib/
git clone https://github.com/afkarxyz/SpotiFLAC.git
mv SpotiFLAC/backend/ lib/
mkdir -p app/ 
cp SpotiFLAC/app.go app/app.go
rm -rf SpotiFLAC

sed -i "s/package main/package app/g" app/app.go
sed -i "s/\"spotiflac\/backend\"/backend \"github.com\/Superredstone\/spotiflac-cli\/lib\"/g" app/app.go
# sed -i "s/backend./lib./g" app/app.go
# sed -i "s/\"github.com\/wailsapp\/wails\/v2\/pkg\/runtime\"//g" $FOLDER/app.go

# for i in $USELESS_FILES; do 
	# rm -rf $FOLDER/$i
# done

for i in $(ls lib/); do
	sed -i "s/package backend/package $FOLDER/g" $FOLDER/$i
done

