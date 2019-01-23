#!/bin/bash

# Adapted from https://github.com/Humpheh/goboy/blob/master/package/darwin/build_darwin.sh

set -e

OUTPUT=../build/SMasteroids.app

# Remove any existing darwin app
if [[ -f ${OUTPUT} ]]; then
    rm -r ${OUTPUT}
fi

# Make the app and the contents directory
mkdir -p ${OUTPUT}

# Copy the contents into the app contents folder
cp -R Contents/ ${OUTPUT}/Contents
sed 's/SHORT_VERSION/'`git describe --abbrev=0 | sed 's/[^0-9\.]*//g'`'/g' Contents/Info.plist | sed 's/VERSION/'`git describe`'/g'  > ${OUTPUT}/Contents/Info.plist

# Build SMasteroids or copy if argument is passed in
if [[ $1 == "" ]]; then
    # Build the binary into the contents executable
    go build -o ${OUTPUT}/Contents/MacOS/SMasteroids gitlab.com/meyerzinn/smasteroids
else
    # Ensure the folder exists
    mkdir -p ${OUTPUT}/Contents/MacOS
    # Move the input file to the contents executable
    cp $1 ${OUTPUT}/Contents/MacOS/SMasteroids
fi

ls ${OUTPUT}/Contents/MacOS/SMasteroids

cp ${OUTPUT}/Contents/MacOS/SMasteroids ../build/SMasteroids
zip -r ../build/SMasteroids.MacOS.App.zip ../build/SMasteroids.app/
zip -r ../build/SMasteroids.MacOS.Binary.zip ../build/SMasteroids

# Cleanup
rm -r ../build/SMasteroids.app/
rm ../build/SMasteroids