# Distributed Systems: Final Project
# Team
|Name|Net ID|
|---|---|
|David Chu|dfc296|
|Smit Sheth|shs572|
|Yash Shah|yjs310|

# Architecture
![Architecture Diagram](https://i.imgur.com/8PTZ8Q5.png)

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

# Running the etcd server
1. Download the [latest etcd binaries](https://github.com/etcd-io/etcd/releases)
2. Install [goreman](https://github.com/mattn/goreman)

    ```go get github.com/mattn/goreman```
3. In the etcd directory, use goreman to start the etcd cluster. Goreman will read the Procfile. 
**Note:-** To run on windows operating system, you need to remove the single quotes (' ') around the cluster definition.

    ```goreman start```

# Project Structure
Project structure is based on [package oriented design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) 
```
yjs310-shs572-dfc296-final-project
│   config.toml     # config file for all services
│   runServers.bat  # batch file to start services
│   
├───cmd
│   ├───authd
│   │       main.go         # start authd server
│   │       
│   ├───frontend\           # start frontend
│   ├───postd
│   │       main.go         # start postd
│   │       
│   ├───userd
│   │       main.go         # start userd
│   │       
│   └───webd
│       │   main.go         # start backend web server
│       │   
│       └───handlers\       # backend web server handlers
│           │   
│           ├───middleware\ # middleware
│           │       
│           └───models\     # request/response models
│                   
├───internal                    # Application logic
│   ├───auth                    # Authentication service
│   │   │   auth_etcd_test.go
│   │   │   auth_test.go
│   │   │   models.go           # Repo Interface definition
│   │   │   
│   │   ├───authentication      # Protobuf files
│   │   │       auth.pb.go
│   │   │       auth.proto
│   │   │       
│   │   ├───service
│   │   │       auth_service.go # Service Implementation
│   │   │       
│   │   └───storage             # Repo Implementations
│   │       ├───etcd
│   │       │       client.go
│   │       │       repository.go
│   │       │       testrepository.go
│   │       │       
│   │       └───memstorage
│   │               memory.go
│   │               repository.go
│   │               testrepository.go
│   │               
│   ├───config                  # Config parsing
│   │       config.go
│   │       
│   ├───post                    # Post service
│   │   │   models.go           # Repo Interface 
│   │   │   post_etcd_test.go
│   │   │   post_test.go
│   │   │   service.go          # Service Implementation
│   │   │   
│   │   ├───postpb
│   │   │       post.pb.go
│   │   │       post.proto
│   │   │       
│   │   └───storage             # Repo Implementations
│   │       ├───etcd
│   │       │       client.go
│   │       │       repository.go
│   │       │       
│   │       └───memstorage
│   │               repository.go
│   │               storage.go
│   │               testrepository.go
│   │               
│   └───user                    # User service
│       │   models.go           # Repo Interface 
│       │   service.go          # Service Implementation
│       │   user_etcd_test.go
│       │   user_test.go
│       │   
│       ├───storage             # Repo Implementations
│       │   │   testrepository.go
│       │   │   
│       │   ├───etcd
│       │   │       client.go
│       │   │       repository.go
│       │   │       
│       │   └───memstorage
│       │           repository.go
│       │           storage.go
│       │           
│       └───userpb              # Protobuf files
│               user.pb.go
│               user.proto
│               
└───vendor
```
