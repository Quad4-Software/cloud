#!/bin/sh
set -e

# init Quad4 Cloud
quad4 init

if [ "$WITH_WRAPPER" = "true" ]; then
    ocwrapper serve --bin=quad4
else
    quad4 server
fi
