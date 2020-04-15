# Distributed Systems: Final Project

# Running the frontend
## Dependencies
- [NodeJS](https://nodejs.org/en/)

- Angular CLI
    `npm install -g @angular/cli`

# Run the project on local system
    cd cmd/frontend
### Install Packages
    npm install

Configure backend url in cmd/frontend/src/environment

### Run
    ng serve

The frontend will listen on port 4200 by default.

# Running the backend
    cd cmd
    runServers.bat

The backend server will listen on port 9090 by default.
The port can be changed in the cmd/webd/config.toml file.
Please note that you have to be in the cmd/webd directory to run the webd.go file because of the way the app searches for the config file

# Project Structure

```
yjs310-shs572-dfc296-final-project
│   config.toml
│   
├───cmd             # Holds files for running the services
│   ├───authd
│   │       main.go # Main file for auth service server
│   │       
│   ├───postd
│   │       main.go # Main file for post service server
│   │       
│   ├───userd
│   │       main.go # Main file for user service server
│   │       
│   └───webd
│       │   main.go # Main file for web server
│       │   
│       └───handlers # web server handlers
│                   
└───internal # Application logic
    ├───auth
    │   ├───authentication      # authentication protobuf files
    │   │       auth.pb.go
    │   │       auth.proto
    │   │       
    │   ├───server              # authentication service
    │   │       auth_server.go
    │   │       
    │   └───storage             # storage for authentication service
    │           memory.go
    │           
    ├───config                  # config loading
    │       config.go
    │       
    ├───post                    
    │   │   models.go
    │   │   service.go          post service 
    │   │   service_test.go
    │   │   
    │   ├───memstorage
    │   │       repository.go
    │   │       storage.go
    │   │       
    │   └───postpb
    │           post.pb.go
    │           post.proto
    │           
    └───user
        │   models.go
        │   service.go
        │   service_test.go
        │   
        ├───memstorage
        │       repository.go
        │       storage.go
        │       
        └───userpb
                user.pb.go
                user.proto
```
