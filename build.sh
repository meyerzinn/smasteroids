packr
go build -ldflags "-X gitlab.com/meyerzinn/smasteroids/smasteroids.version=$(git describe)" -o build/smasteroids gitlab.com/meyerzinn/smasteroids