#!/usr/bin/env bash

OUTPUT=../build

cp $1 ${OUTPUT}/SMasteroids.exe
zip -r ../build/SMasteroids.Windows.zip ${OUTPUT}/SMasteroids.exe
rm ${OUTPUT}/SMasteroids.exe