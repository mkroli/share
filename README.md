share
=====

share makes local files available via http.

## Installation

```sh
go get github.com/mkroli/share
```

## Usage

### Server mode
The first time a users starts share it'll run in server mode.
This means it starts the http server which will share all files given in the parameter list.


### Client mode
The following times the user starts share it'll communicate all files given in the parameter list to the server instance which will then share those files as well.
