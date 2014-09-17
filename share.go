/*
 * Copyright 2014 Michael Krolikowski
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printError() {
	if err := recover(); err != nil {
		fmt.Printf("\033[0;31m\033[1mError\033[0m: %s\n", err)
	}
}

func printLocation(filename, location string) {
	fmt.Printf("%s is shared at \033[1m%s\033[0m\n", filename, location)
}

var nextFileNumber = 0
var files = make(map[string]string)
var filesMutex sync.Mutex

func addFile(filename string) string {
	filesMutex.Lock()
	defer filesMutex.Unlock()
	nextFileNumber++
	files[path.Join(strconv.Itoa(nextFileNumber), filepath.Base(filename))] = filename
	return fmt.Sprintf("http://%s:%s/%d/%s", host, port, nextFileNumber, filepath.Base(filename))
}

var host, port string

func main() {
	defer printError()
	var index bool
	hostname, err := os.Hostname()
	panicOnError(err)
	flag.Usage = func() {
		fmt.Printf("Usage: %s [file]...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&host, "host", hostname, "the host to bind to")
	flag.StringVar(&port, "port", "8080", "the port to bind to")
	flag.BoolVar(&index, "index", false, "show list of all shared files")
	flag.Parse()

	username, err := user.Current()
	panicOnError(err)
	addr, err := net.ResolveUnixAddr("unix", filepath.Join(os.TempDir(), "share-"+username.Username+".sock"))
	panicOnError(err)
	l, err := net.ListenUnix("unix", addr)
	if err == nil {
		defer l.Close()
		server(l, port, index)
	} else {
		client(addr)
	}
}
