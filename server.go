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
	"bufio"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

func handleIpcConnection(conn net.Conn) {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		filename := strings.TrimSuffix(line, "\n")
		location := addFile(filename)
		printLocation(filename, location)
		_, err = w.WriteString(location + "\n")
		if err != nil {
			break
		}
		err = w.Flush()
		if err != nil {
			break
		}
	}
}

func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	filesMutex.Lock()
	defer filesMutex.Unlock()

	filename, ok := files[r.URL.Path[1:]]
	if ok {
		http.ServeFile(w, r, filename)
	} else {
		http.NotFound(w, r)
	}
}

func server(l *net.UnixListener, httpPort string) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	go func() {
		<-signalChan
		l.Close()
		os.Exit(0)
	}()

	for _, arg := range flag.Args() {
		filename, err := filepath.Abs(arg)
		panicOnError(err)
		location := addFile(filename)
		printLocation(filename, location)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				go handleIpcConnection(conn)
			}
		}
	}()

	http.HandleFunc("/", handleHttpRequest)
	err := http.ListenAndServe(":"+httpPort, nil)
	panicOnError(err)
}
