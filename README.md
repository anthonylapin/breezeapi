# BreezeAPI

## Description

BreezeAPI is a lightweight HTTP framework that was built as a small pet project aimed to build an HTTP server that's capable of handling simple GET/POST requests, serving files, handling concurrent connections and compressing responses using GZIP

Features:

1. Extracting URL path and dynamic params wrapped within {paramName}
2. Reading headers/body of the request
3. Support for concurrent connections using workers pool with Go channels
4. Returning files
5. Gzip compression

Framework is built in pretty easy way similar to already existing frameworks such as express in Node.js:

1. Create server object
2. Create a router for the server
3. Add endpoint handlers using router.Get/Post/{otherMethod}() methods.
4. Add router to the server
5. Start listening from the server

## How to run locally:

You need to have a golang installed on our machine compatible with the version specified in go.mod.

1. Clone repository
2. Run using start_server.sh script
3. Framework code is made within server directory
4. main.go contains testing endpoints
