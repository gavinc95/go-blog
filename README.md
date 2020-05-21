## Overview
This is a simple blogging application built with an HTTP server and PostgreSQL.

The HTTP server handles requests using a [Mux Router](https://github.com/gorilla/mux).
It uses the [lib/pq](https://github.com/lib/pq) driver to interact with the database.
PostgreSQL is run within a Docker container, for easy end to end integration tests.

## Requirements
To use this application, you need to have the following installed:
- PostgreSQL 
	- Environment variables:
	```
	export POSTGRES_USER=postgres
	export POSTGRES_PASSWORD=password
	export APP_DB_NAME=postgres
	```
	- If you are using MacOS, you can follow instructions [here](https://gist.github.com/ibraheem4/ce5ccd3e4d7a65589ce84f2a3b7c23a3) to install and run PostgreSQL.
- Docker
- Go
	- Go Modules and Private repo setup (may not need this):
	```
	export GO111MODULE=on
	export GOPROXY=direct
	export GOSUMDB=off
    export GIT_TERMINAL=1
	```

## Usage
First, we want to run the PostgreSQL instance in a Docker image like so:
```
docker run -it -p 5432:5432 -d postgres
```

Then, in the project root, run `go build` and `./go-blog` to start the application binary.

In a separate terminal, you can then make relevant API calls.
The requests can be made through cURL commands to the `/users` or the `/posts` endpoints using different requests in JSON format.

For example:
A `GET` request to the `/users` endpoint can send a `GetUserRequest` to the server like so:
```
curl -X GET localhost:8080/users -d '{"id": "<USER_ID>"}'
```
or you can send `POST` with a `CreateUserRequest`
```
curl -X POST localhost:8080/users -d '{"name": "<NAME>", "email": "<EMAIL>"}'
```

### Requests
This app only supports CRUD operations for a blog via `User` and `Post` [models](https://github.com/gavinc95/go-blog/blob/master/db/models/models.go).

The requests are defined below:
```
type GetUserRequest struct {
	ID string `json:"id"` // required
}

type GetUserResponse struct {
	User *models.User `json:"user"`
}

type CreateUserRequest struct {
	Email string `json:"email"` // required
	Name  string `json:"name"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type UpdateUserRequest struct {
	ID    string `json:"id"` // required
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateUserResponse struct {
	ID string `json:"id"`
}

type DeleteUserRequest struct {
	ID string `json:"id"` // required
}

type DeleteUserResponse struct {
	ID string `json:"id"`
}

type CreatePostRequest struct {
	UserID  string `json:"user_id"` // required
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreatePostResponse struct {
	ID string `json:"id"`
}

type UpdatePostRequest struct {
	ID      string `json:"id"` // required
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostResponse struct {
	ID string `json:"id"`
}

type GetPostRequest struct {
	ID string `json:"id"` // required
}

type GetPostResponse struct {
	Post *models.Post `json:"post"`
}

type GetAllPostsRequest struct {
	UserID string `json:"user_id"` // required
}

type GetAllPostsResponse struct {
	Posts []*models.Post `json:"posts"`
}

type DeletePostRequest struct {
	ID string `json:"id"` // required
}

type DeletePostResponse struct {
	ID string `json:"id"`
}
```

## Notes
This is far from a complete blog management platform. Some notable things that weren't addressed are: 
- Authentication - we could store encrypted passwords or use a JWT for each user, and validate during a login step.
- The post content doesn't support images or audio, but if we did, we could store them with the following schema:
	- `imageID` -> `S3 URI`, and actually store the image in an object store like Amazon S3 or Google Cloud Storage.
- Access Control - e.g. even if one user knows the `postID` that belongs to another user, they shouldn't be able to edit/delete that content.
- Deployment
- Anything scalability related, like load balancing requests, etc.
