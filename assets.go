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

var indexTpl = `<!DOCTYPE html>
<html>
<head>
<title>share - Index</title>
<style type="text/css">
h1 {
	border-bottom: 1px solid grey;
}

li {
	display: block;
}

a {
	color: rgba(0, 0, 0, .75);
	text-decoration: none;
}

a:hover {
	color: rgba(0, 0, 0, 1);
	text-decoration: underline;
}
</style>
</head>
<body>
	<h1>shared files</h1>
	<ul>{{range .}}
		<li><a href="{{.}}">{{.}}</a></li>{{end}}
	</ul>
</body>
</html>`
