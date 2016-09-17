#!/bin/sh -e
go get -u golang.org/x/net/context
go get -u google.golang.org/api/container/v1
go get -u k8s.io/client-go/1.4/kubernetes
go get -u google.golang.org/grpc
go get -u golang.org/x/oauth2
go get -u cloud.google.com/go/compute/metadata