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
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
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

var indexTemplate, _ = template.New("index").Parse(indexTpl)

func handleIndexHttpRequest(w http.ResponseWriter, r *http.Request) {
	filesMutex.Lock()
	defer filesMutex.Unlock()

	links := make([]string, len(files))
	i := 0
	for l, _ := range files {
		links[i] = l
		i++
	}
	sort.Strings(links)
	indexTemplate.Execute(w, links)
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

func handlerFactory(index bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && index {
			handleIndexHttpRequest(w, r)
		} else {
			handleHttpRequest(w, r)
		}
	}
}

func server(l *net.UnixListener, httpPort string, index bool) {
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

	http.HandleFunc("/", handlerFactory(index))
	err := http.ListenAndServe(":"+httpPort, nil)
	panicOnError(err)
}
