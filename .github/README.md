# CAM
Cold's Achievement Manager

A synopsis: [SAM](https://github.com/gibbed/SteamAchievementManager) is cool, but there are some small issues with it;
- It doesn't really support automation.
- Windows only

And these are the two things I wanted to fix.

# How to use

### Method A: Use Go to run it
You'll need
- An internet connection
- A web browser
- Go

Download [Go](https://go.dev/dl) and follow the [installation instructions](https://go.dev/doc/install)
then, clone (`git clone https://github.com/colduw/cam`) (or [download](https://github.com/colduw/cam/archive/refs/heads/main.zip)) the repository.
Afterwards, it's a few commands:
```bash
cd cam
go run main.go --appID=480
```
And you're done. For any other information, run
```bash
go run main.go --help
```

### Method B: Pre-built binaries
This section is a WIP

# Attribution
- https://github.com/ebitengine/purego - Used for using the SteamAPI's exported functions - Licensed under the [Apache Version 2.0 License](https://github.com/ebitengine/purego/blob/main/LICENSE)
- https://github.com/charmbracelet/log - Used for logging - Licensed under the [MIT License](https://github.com/charmbracelet/log/blob/main/LICENSE)
- https://github.com/charmbracelet/lipgloss - Used for logging - Licensed under the [MIT License](https://github.com/charmbracelet/lipgloss/blob/master/LICENSE)