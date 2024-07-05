# User REST API

This User API is implemented in Go. It can store and retrieve a user's name, age, and email. It implements the following methods:


POST /api/v1/user  
Requires HTTP Basic Auth (username and password strings are extracted from the header)  
The body must be a json with fields: "age" int, "name" string, "email" string  

GET /api/v1/user/:username  
Requires HTTP Basic Auth  


PUT /api/v1/user/:username  
Requires HTTP Basic Auth  
The body must be a json with fields: "age" int, "name" string, "email" string  

DELETE /api/v1/user/:username  
Requires HTTP Basic Auth  

This API meets the requirements of a REST API.  
- All methods are stateless (do not depend on previous requests)
- Repeated GET requests always return same resource (can be cached)
- Repeated PUT and DELETE requests have the same effect as one request

## Steps to run
Install Go

Install Docker

Set environment variable $MYSQL_ROOT_PASSWORD

If on Linux run the shell script to start mysql docker container
If on Windows copy the shell script and run on command line terminal

In root directory, run `go run ./server`

To test, run `go test ./server ./database`