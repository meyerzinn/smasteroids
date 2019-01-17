packr
go build -ldflags "-X gitlab.com/meyerzinn/smasteroids/smasteroids.version=$(git describe)" -o mac/SMasteroids.app/Contents/MacOS/smasteroids gitlab.com/meyerzinn/smasteroids
