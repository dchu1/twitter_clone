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

# Running the backend
    go run cmd/webd/webd.go

The backend server will listen on port 9090 by default.
The port can be changed in the cmd/webd/config.toml file.

# Project Structure

```
github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project
├───cmd             # Application code
│   │   webd.go         # Main file that runs the backend server
│   │   config.toml     # Configuration file
│   ├───frontend        # Frontend Application code
│   └───webd            # Backend Application code
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