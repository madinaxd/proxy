# Simple proxy server

A server receives a request /curl from client. Request body should be in JSON format. Example: 

```console
{
    "method": "GET",
    "url": "http://google.com",
    "headers": {
        "Authentication": "Basic bG9naW46cGFzc3dvcmQ=",
        ....
    }
}
```

Server forms valid HTTP-request to 3rd-party service with data from client's message and responses to client with JSON object:

```console
{
    "id": <generated unique id>,
    "status": <HTTP status of 3rd-party service response>,
    "headers": {
        <headers array from 3rd-party service response>
    },
    "length": <content length of 3rd-party service response>
}
```


## Local Development

Run the commands below.

```console
go run main.go
```

## Test
Now you can test the API with Postman post method at:

```console
localhost:3000
```