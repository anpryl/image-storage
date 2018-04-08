# image-storage
Image storing solution  
Uses [traefik](https://traefik.io/) as reverse proxy  
Consist of two services:
* [authsvc](https://github.com/anpryl/image-storage/tree/master/authsvc) available at /authsvc
* [imagesvc](https://github.com/anpryl/image-storage/tree/master/imagesvc) available at /imagesvc

## Requirements
[Docker](https://docs.docker.com/install)
[docker-compose](https://docs.docker.com/compose/install)
[httpie](https://httpie.org/)(for testing, curl has issues with file upload)

## Build and deploy
```
docker-compose up
```

## Usage example

```
# Check jwt middleware
$ http GET http://127.0.0.1/imagesvc/images Authorization:1
HTTP/1.1 401 Unauthorized
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 23:04:24 GMT

# Obtain token
$ http http://127.0.0.1/authsvc/token
HTTP/1.1 200 OK
Content-Length: 118
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:46:21 GMT

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg"
}

# Check if storage is empty
$ http GET http://127.0.0.1/imagesvc/images Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg
HTTP/1.1 200 OK
Content-Length: 386
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:49:29 GMT

{
    "images": null
}

# Upload image
$ http POST http://127.0.0.1/imagesvc/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg < imagesvc/Simon-Peyton-Jones-feature.jpg
HTTP/1.1 201 Created
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:52:19 GMT

# Check list of all images again
$ http GET http://127.0.0.1/imagesvc/images Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg
HTTP/1.1 200 OK
Content-Length: 17
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:54:55 GMT

{
    "images": [
        "1"
    ]
}

# Check if it exist
$ http GET http://127.0.0.1/imagesvc/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg
HTTP/1.1 200 OK
Content-Type: image/jpeg
Date: Sun, 08 Apr 2018 22:53:14 GMT
Transfer-Encoding: chunked

+-----------------------------------------+
| NOTE: binary data not shown in terminal |
+-----------------------------------------+

# Download image to /tmp/image1
$ http GET http://127.0.0.1/imagesvc/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg > /tmp/image1

# Deleting image
$ http DELETE http://127.0.0.1/imagesvc/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg
HTTP/1.1 200 OK
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:56:00 GMT

# Check list of all images again
$ http GET http://127.0.0.1/imagesvc/images Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg
HTTP/1.1 200 OK
Content-Length: 16
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:56:33 GMT

{
    "images": null
}

# Try to upload
$ http GET http://127.0.0.1/imagesvc/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMjMyMjc1ODF9.5-7Q-7wyANVsJgtgi-TX4FsUyEy-v2Bzph0L0FBKuKg
HTTP/1.1 404 Not Found
Content-Length: 40
Content-Type: text/plain; charset=utf-8
Date: Sun, 08 Apr 2018 22:56:52 GMT

{
    "message": "File not found"
}
```
