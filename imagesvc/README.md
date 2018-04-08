# imagesvc
Service storing and serving images for image-storage

## Requirements
[Docker](https://docs.docker.com/install)
[httpie](https://httpie.org/)(for testing, curl has issues with file upload)

## Build
```
./build.sh
```

## Usage
Available options:  
`docker run imagesvc --help`  
Start:  
`docker run -p 4000:4000 imagesvc -secret super_secret_phrase -host 0.0.0.0`

## API

Every endpoint check ~Authorization~ header, it should contain signed token

### Errors

* Empty filename
* File exist
* File not found

### Save image
POST /images/:name  
Response:  
Image saved (status 201) or error(details in body)

Example:  
```
$ http http://127.0.0.1:4000/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyMjY2MzF9.bHWykLyOeDVOW2VWeWNhSeySDFUh8jL0KeYuHkvf2YM < Simon-Peyton-Jones-feature.jpg
```

### Images names
GET /images  
Response:  
```json
{
  "images": [
    "1"
  ]
}
```

Example:  
```
$ http http://127.0.0.1:4000/images Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyMjY2MzF9.bHWykLyOeDVOW2VWeWNhSeySDFUh8jL0KeYuHkvf2YM
{
    "images": [
        "1",
    ]
}
```

### Image by name
GET /images/:name  
Response:  
Image in body or 404

Example:  
```
$ http http://127.0.0.1:4000/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyMjY2MzF9.bHWykLyOeDVOW2VWeWNhSeySDFUh8jL0KeYuHkvf2YM > /tmp/testimage
```

### Delete image
DELETE /images/:name  
Response:  
Image deleted (status 200) or error(details in body)

Example:  
```
$ http DELETE http://127.0.0.1:4000/images/1 Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMyMjY2MzF9.bHWykLyOeDVOW2VWeWNhSeySDFUh8jL0KeYuHkvf2YM
```
