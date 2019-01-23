packr
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-X github.com/20zinnm/smasteroids/smasteroids.version=$(git describe)" -o dist/macos/SMasteroids.app/Contents/MacOS/smasteroids github.com/20zinnm/smasteroids
CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-X github.com/20zinnm/smasteroids/smasteroids.version=$(git describe)" -o dist/windows/smasteroids.exe github.com/20zinnm/smasteroids

