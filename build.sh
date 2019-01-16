packr build -ldflags "-X gitlab.com/meyerzinn/smasteroids/game.version=$(git describe)" -o build/smasteroids gitlab.com/meyerzinn/smasteroids
