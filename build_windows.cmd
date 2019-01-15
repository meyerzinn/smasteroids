packr build -ldflags "-X gitlab.com/meyerzinn/smasteroids/game.version=$(git describe --abbrev=0 --tags)" -o build/smasteroids.exe gitlab.com/meyerzinn/smasteroids
