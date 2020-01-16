# http-monitor

[![Build Status](https://cloud.drone.io/api/badges/smf8/http-monitor/status.svg)](https://cloud.drone.io/smf8/http-monitor)

A HTTP endpoint monitor service written in go with RESTful API.

ORM library [Gorm](https://github.com/jinzhu/gorm)

Web framework [Echo](https://echo.labstack.com/)

Job queue manager [workerpool](https://github.com/gammazero/workerpool)

struct and data validator [Govalidator](https://github.com/asaskevich/govalidator)

- [Installation](#Installation)
- [Database](#Database)
- [API](#API)
- [Package Structure](#Package-Structure)

## Installation

- Make sure go is installed properly.

- Download project into your GOPATH

  ```
  $ go get -u -v github.com/smf8/http-monitor
  $ cd $GOPATH/src/github.com/smf8/http-monitor
  ```
  
- Build project using `go build main.go`

- You can also run the project using docker-compose by running `docker-compose up` ( change container port in `docker-compose.yml`)

## Database

#### Tables : 

**Users:**

| id(pk)  | created_at | updated_at | deleted_at | username     | password     |
| :------ | ---------- | ---------- | ---------- | ------------ | ------------ |
| integer | datetime   | datetime   | datetime   | varchar(255) | varchar(255) |

**URLs:**

| id(pk)  | created_at | updated_at | deleted_at | user_id(fk) | address      | threshold | failed_times |
| ------- | ---------- | ---------- | ---------- | ----------- | ------------ | --------- | :----------- |
| integer | datetime   | datetime   | datetime   | integer     | varchar(255) | integer   | integer      |

**Requests:**

| id(pk)  | created_at | updated_at | deleted_at | url_id(fk) | result  |
| ------- | ---------- | ---------- | ---------- | ---------- | ------- |
| integer | datetime   | datetime   | datetime   | integer    | integer |

## API

### Specs:

For all requests and responses we have `Content-Type: application/json`.

Authorization is with JWT.

#### User endpoints:

**Login:**

`POST /api/users/login`

request structure: 

```
{
	"username":"foo", // alpha numeric, length >= 4
	"password":"*bar*" // text, length >=4 
}
```

**Sign Up:**

`POST /api/users`

request structure (same as login):

```
{
	"username":"foo", // alpha numeric, length >= 4
	"password":"*bar*" // text, length >=4 
}
```

#### URL endpoints:

**Create URL:**

`POST /api/urls`

request structure:

```
{
	"address":"http://some-valid-url.com" // valid url address
	"threshold":20 // url fail threshold
}
```

##### **Get user URLs:**

`GET /api/urls`

**Get URL stats:**

`GET /api/urls/:urlID?from_time&to_time`

`urlID` a valid url id

`from_time` a starting time in unix time format(Optional, `to_time` only is not allowed) 

`to_time` an ending time in unix time format.(Optoinal)

**Delete URL:**

`DELETE /api/urls/:urlID`

`urlID` a valid url id to be deleted

**Get URL alerts:**

`GET /api/alerts`

**Dismiss URL alerts:**

`PUT /api/alerts/:urlID`

`urlID` a valid url. **This endpoint reset given url's failed_times to 0 ** 

#### Responses:

##### Errors:

If there was an error during processing the request, a json response with the following format is returned with related response code: 

```
{
	"errors":{
		"key":"value" // a list of key,value of errors occurred
	}
}
```

##### URL stat:

```
{
    "data": {
        "url": "http://google.com",
        "requests_count": 1,
        "requests": [
            {
                "result_code": 200,
                "created_at": "2019-01-16T14:07:25.443300581+03:30"
            }
        ]
    }
}
```

##### List of URLs:

```
{
    "data": {
    	"url_count": 1,
        "urls": [
            {
                "id": 0,
                "url": "http://google.com",
                "user_id": 1,
                "created_at": "2020-01-16T14:07:15.066047519+03:30",
                "threshold": 10,
                "failed_times": 0
            }
        ]
    }
}
```

##### Request report:

```
{
	"data": "A message with report"
}
```



## Package Structure

```
├── common				// common package for commonly used functions
│   ├── erros.go
│   └── jwt.go
├── db					// database creation and initialization
│   └── db.go
├── handler				// handler package for routing and request handling
│   ├── handler.go
│   ├── request.go
│   ├── response.go
│   ├── routes.go
│   ├── url.go
│   └── user.go
├── main.go				// main entry of application
├── middleware			// middlewares used in API
│   └── jwt.go
├── model				// data types used in application
│   ├── model_test.go
│   ├── url.go
│   └── user.go
├── monitor				// monitor package to handle url monitoring and scheduling
│   ├── monitor.go
│   ├── monitor_test.go
│   └── scheduler.go
└── store				// a layer for model-database interactions
    ├── store.go
    └── store_test.go
```

#### TODO

- [ ] Refactor error management system
- [ ] Improve scheduler to accept different intervals
- [ ] Integrate project with a configuration library 