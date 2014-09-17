share
=====

share makes local files available via http.

## Installation

Either [Download binary (Linux x86_64)](https://github.com/mkroli/share/releases/download/0.1/share)
or build it:
```sh
go get github.com/mkroli/share
```

## Usage
```
Usage: share [file]...
  -host="host": the host to bind to
  -index=false: show list of all shared files
  -port="8080": the port to bind to
```

### Server mode
The first time a users starts share it'll run in server mode.
This means it starts the http server which will share all files given in the parameter list.


### Client mode
The following times the user starts share it'll communicate all files given in the parameter list to the server instance which will then share those files as well.
