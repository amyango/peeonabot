# peeonabot

## Requirements
- Docker  
This is how I did it for OSX:
```
brew install golang
brew install docker
brew install colima
colima start
```

## Run instructions
1. Get the bot API key from the Discord Developer GUI and 
   copy it into `credentials/discord.token`
2. Build by running `./reload`
3. You can check the log output but running:
```
docker logs peeonabot
```
Use the -f flag for following log output (shows new log output
as it comes in)

