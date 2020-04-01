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
    cd cmd/webd
    go run webd.go

The backend server will listen on port 9090 by default.
The port can be changed in the cmd/webd/config.toml file.
Please note that you have to be in the cmd/webd directory to run the webd.go file because of the way the app searches for the config file

# Project Structure

```
github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project
├───cmd             # Application code
│   ├───frontend        # Frontend Application code
│   └───webd            # Backend Application code
│       │   webd.go         # Main file that runs the backend server
│       │   config.toml     # Configuration file
│       ├───app             # Package for manipulating our data structures.
│       ├───auth            # Package for managing sessions
│       │   ├───session
│       │   └───storage
│       ├───config          # Package for handling configuration files
│       ├───handlers        # Package with http handlers
│       │   ├───middleware      # Package with middleware code
│       │   └───models          # Package with request and response struct definitions       
├───internal        # Packages that are shared across applications
└───vendor          # Vendored 3rd Party dependencies
```