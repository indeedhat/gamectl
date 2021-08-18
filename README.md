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
- [ ] server logs
    - [ ] add to config
        - should be able to have multple log files/commands
    - [ ] add websocket endpoint for streaming logs
    - [ ] add option to frontend (show in modal)
- [ ] maybe add some unit testing
- [ ] docker file
    - this may not be possible given that it will need to run commands on the host machine
- [ ] docker-compose config
- [ ] reload app config

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

