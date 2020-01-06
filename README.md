# http-monitor

[![Build Status](https://cloud.drone.io/api/badges/smf8/http-monitor/status.svg)](https://cloud.drone.io/smf8/http-monitor)

A HTTP endpoint monitor service written in go with RESTful API.

### Database:

Database ORM library: [gorm](https://github.com/jinzhu/gorm)

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

### Structure:

```
├── db 					// database layer
│   ├── db.go
│   └── db_test.go
├── model 				// models
│   ├── model_test.go
│   ├── url.go
│   └── user.go
└── store				// layer for models and database interaction
│   └── store.go
├── go.mod
├── go.sum
├── main.go				//main entry of application
```