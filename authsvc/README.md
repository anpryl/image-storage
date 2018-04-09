# authsvc
Service generate access tokens for image-storage

## Requirements
[Docker 17.09+](https://docs.docker.com/install)  

## Build
```
./build.sh
```

## Usage
Available options:  
`docker run authsvc --help`  
Start:  
`docker run -p 4000:4000 authsvc -secret super_secret_phrase -host 0.0.0.0`

## API

### Generate new token:  
GET /token  
Response:  
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjMyMTg5NTZ9.wHEYwK0XUUSHi7mAT4Q0ZD0Mr5trs1oAcTaCsykdyfM"
}
```

Example:  
```
$ curl -X GET http://127.0.0.1:4000/token
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjMyMTk2NTJ9.cAVvNdRwJyUvMO0DvH7K-v0iFDkqq18VahObV-wj9EE"}
```
