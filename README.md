# Command Center
Control game servers from web ui

## Requirements
go 1.16+ \
cmake

## Run the project
```shell
make run
```

## TODO
- [ ] implement start/stop restart in frontend
- [ ] implement status in frontend
- [ ] game server requirements config
    - i want to allow game server configs to be able to define the resources
      they require, that way cc can make decisions on if multiple can run at once
- [x] activity monitor
- [x] server logs
    - [x] add to config
    - [x] should be able to have multple log files/commands
    - [x] add websocket endpoint for streaming logs
        - ended up using http streams instead of web sockets
    - [x] add option to frontend (show in modal)
- [x] reload app config
- [ ] maybe add some unit testing
- [ ] docker support
    - running this in docker would be a pain as it would require a companion app to run on the host to trigger\
      the scripts, it is possible and i may look into it
    - [ ] companion app
    - [ ] implement into web app
    - [ ] docker file
    - [ ] docker-compose config

The webui will  be available at localhost:8080

## Kagorus jobs
Design a pretty user interface

for now this can all be done in the index.tmpl file

It would be good to get you to do some of the intergrating the theme into the system but as
you have never actually used go before (and more to the point its templating engine) i might need
to go through that with you to begin with.

## Directory structure
- **app/** this is where the backend logic exists
- **assets/** this is where all the public available assets exist (css, js, images etc.)
- **config/** this is where the yaml files for setting up the game servers live
- **views/** this is where the template files exist (they need to be .tmpl files currently but i can change that to .html if you prefer)
- **.env** the contents of this file get loaded as environment variables on boot

