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
    ./runServers.bat

The backend server will listen on port 9090 by default.
The port can be changed in the config.toml file.
Please note that if you want to use the config.toml file, you need to start services from the root folder (yjs310-shs572-dfc296-final-project/). This is because the services look in the current working directory
for the config file. Otherwise they will use default values.

# Project Structure
Project structure is based on [package oriented design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) 
```
yjs310-shs572-dfc296-final-project
│   config.toml     # config file for all services
│   runServers.bat  # batch file to start services
│   
├───cmd             # Holds files for starting the services
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
    │   ├───servic              # authentication service
    │   │       auth_service.go
    │   │       
    │   └───storage             # storage for authentication service
    │           memory.go
    │           
    ├───config                  # config loading
    │       config.go
    │       
    ├───post                    
    │   │   models.go
    │   │   service.go          # post service 
    │   │   
    │   ├───memstorage          # storage for post service
    │   │       repository.go
    │   │       storage.go
    │   │       
    │   └───postpb              # protobuf files for post service
    │           post.pb.go
    │           post.proto
    │           
    └───user                    # user service
        │   models.go            
        │   service.go
        │   user_test.go
        │   
        ├───memstorage          # storage for user service
        │       repository.go
        │       storage.go
        │       
        └───userpb              # protobuf files for user service
                user.pb.go
                user.proto
```
