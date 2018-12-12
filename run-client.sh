#!/bin/bash

rm -f ../../bin/client
go install grpc-demo/client
../../bin/client -o $1
