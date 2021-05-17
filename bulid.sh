
# pack resources
# gf pack config,public,template packed/data.go -n packed
gf pack public,template packed/data.go -n packed

# pack swagger
gf swagger --pack

# build application
# gf build main.go -n locksim -v 0.1.1 -a amd64,386 -s linux,windows,darwin -p ./build
gf build main.go -n locksim -v 0.1.1 -a amd64 -s linux,windows -p ./build