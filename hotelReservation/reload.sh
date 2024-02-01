go mod vendor

rm -rf vendor/github.com/AleckDarcy/reload

mkdir -p vendor/github.com/AleckDarcy/

ln -s $GOPATH/src/github.com/AleckDarcy/reload vendor/github.com/AleckDarcy/reload