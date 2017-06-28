#!/bin/sh

cd $APPDIR
chmod +x main
./main -contest=${CONTEST} -pushHost=${PUSHHOST} -port=${PORT}