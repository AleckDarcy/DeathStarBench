CBPATH="/Users/aleck/go/src/github.com/AleckDarcy/reload"

docker compose down

echo process vendor
rm -rf vendor/github.com/AleckDarcy/reload
mkdir -p vendor/github.com/AleckDarcy/
ln -s $CBPATH vendor/github.com/AleckDarcy/reload