CBPATH="/Users/aleck/go/src/github.com/AleckDarcy/ContextBus"

docker compose down

echo process vendor
rm -rf vendor/github.com/AleckDarcy/ContextBus
cp -r $CBPATH vendor/github.com/AleckDarcy/ContextBus

echo compling and start
docker compose up -d --build

echo process vendor
rm -rf vendor/github.com/AleckDarcy/ContextBus
mkdir -p vendor/github.com/AleckDarcy/
ln -s $CBPATH vendor/github.com/AleckDarcy/ContextBus