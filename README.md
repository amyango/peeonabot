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
   copy it into `credentials/peeonabot-prod.token`
   If you want a test version to add to a test server, also
   copy the test bot's credentials into 
   `credentials/peeonabot-test.token`.
2. Build by running `./reload prod` (or `./reload test` for
   a test bot)
3. You can check the log output but running:
```
docker logs peeonabot-prod
docker logs peeonabot-test
```
Use the -f flag for following log output (shows new log output
as it comes in)

