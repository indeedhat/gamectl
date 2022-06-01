# GameCTL
Control game servers from web ui

## Requirements
go 1.15+ `may work with earlier version but untested` \
cmake

## Run the project
```shell
make build run
```
The webui will be available at localhost:8080

## Issues
- [ ] long polling for status sometimes breaks with a 500 error

## TODO
- [ ] game server requirements config
    - i want to allow game server configs to be able to define the resources
      they require, that way gctl can make decisions on if multiple can run at once
- [ ] interactive server shell
- [ ] new server wizzard
- [x] Cron tasks
    - [x] restart server
    - [x] backup server
- [x] implement start/stop restart in frontend
- [x] implement status in frontend
- [x] world backups
- [x] activity monitor
- [x] server logs
    - [x] add to config
    - [x] should be able to have multple log files/commands
    - [x] add websocket endpoint for streaming logs
        - ended up using http streams instead of web sockets
    - [x] add option to frontend (show in modal)
- [x] reload app config


## Directory structure
- **.temp/** used for creating temp files during world download process  
- **.config/** this is where the yaml files for setting up the game servers live  
- **cmd/** this is where the main function for the binary lives
- **configs/** example crap for congfigs  
- **internal/** this is where the backend logic exists  
- **web/assets/** this is where all the public available assets exist (css, js, images etc.)  
- **web/views/** this is where the template files exist (they need to be .tmpl files currently but i can change that to .html if you prefer)  
- **.env** the contents of this file get loaded as environment variables on boot 

