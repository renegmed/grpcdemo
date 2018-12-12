#!/bin/bash

rm -f ../../bin/server
go install grpc-demo/server
../../bin/server