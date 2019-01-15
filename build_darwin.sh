packr build -ldflags "-X gitlab.com/meyerzinn/smasteroids/game.version=$(git describe --abbrev=0 --tags)" -o build/smasteroids gitlab.com/meyerzinn/smasteroids
