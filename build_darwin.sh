packr build -ldflags "-X gitlab.com/meyerzinn/smasteroids/game.version=$(git describe --abbrev=0 --tags)" gitlab.com/meyerzinn/smasteroids
