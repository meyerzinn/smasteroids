packr
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-X gitlab.com/meyerzinn/smasteroids/smasteroids.version=$(git describe)" -o dist/darwin/SMasteroids.app/Contents/MacOS/smasteroids gitlab.com/meyerzinn/smasteroids
CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-X gitlab.com/meyerzinn/smasteroids/smasteroids.version=$(git describe)" -o dist/windows/smasteroids.exe gitlab.com/meyerzinn/smasteroids

