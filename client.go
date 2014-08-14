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
	"path/filepath"
	"strings"
)

func client(addr *net.UnixAddr) {
	conn, err := net.DialUnix("unix", nil, addr)
	panicOnError(err)
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	for _, arg := range flag.Args() {
		filename, err := filepath.Abs(arg)
		panicOnError(err)
		_, err = w.WriteString(filename + "\n")
		panicOnError(err)
		err = w.Flush()
		panicOnError(err)
		line, err := r.ReadString('\n')
		printLocation(filename, strings.TrimSuffix(line, "\n"))
	}
	conn.Close()
}
