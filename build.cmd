rem this build script is just for local testing, there is an automatic build that happens in github
rem and good to be aware: these version numbers get populated/incremented in that github build
go-winres make --file-version=1.0.1.0 --product-version=1.0.1.0
go build -ldflags "-X main.version=1.0.1" || pause
