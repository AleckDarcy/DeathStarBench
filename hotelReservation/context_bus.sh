go mod vendor

rm -rf vendor/github.com/AleckDarcy/ContextBus

mkdir -p vendor/github.com/AleckDarcy/

ln -s $GOPATH/src/github.com/AleckDarcy/ContextBus vendor/github.com/AleckDarcy/ContextBus