## CronHubot

* All the flexibility and power of Cron as a Service.
* Simple REST protocol, integrating with a web application in a easy and straightforward way.
* No more wasting time building and managing scheduling infrastructure.

## Basic Concepts
Works by calling back to your application via HTTP GET according to a schedule constructed by you or your application.
Dynamic register jobs, low-footprint, durable.

## Setup
Env vars
```bash
// postgresql
export DATASTORE_URL="postgresql://postgres@localhost/dbname?sslmode=disable"
export SERVICE_PORT=3000

// mysql
export DATASTORE_URL="mysql://root:123456@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
export SERVICE_PORT=3000
```

## Authentication
This API does not ship with an authentication layer. You **should not** expose the API to the Internet. This API should be deployed behind a firewall, only your application servers should be allowed to send requests to the API.

## API Endpoints
- [`GET` /health](#get-health) - Get application health
- [`GET` /cronjobs](#get-cronjobs) - Get a list of scheduled cronjobs
- [`POST` /cronjobs](#post-cronjobs) - Create a cronjob
- [`GET` /cronjobs/:id](#get-cronjobsid) - Get a single cronjob
- [`DELETE` /cronjobs/:id](#delete-cronjobsid) - Delete a cronjob
- [`PATCH` /cronjobs/:id](#patch-cronjobsid) - Update a cronjob

### API Documentation
#### `GET` `/cronjobs`
Get a list of available cronjobs.
- Method: `GET`
- Endpoint: `/cronjobs`
- Responses:
    * 200 OK
    ```json
    [
       {
          "id":1,
          "name":"cronjob-1",
          "url":"your-api/job",
          "expression": "0 5 * * * *",
          "status": 1,
          "retries": 2,
          "request_timeout": 3,
          "stop_at": "2019-11-11 01:00:00",
          "created_at": "2016-12-10 14:02:37",
          "updated_at": "2016-12-10 14:02:37"
       }
    ]
    ```
    - `id`: is the id of the cronjob.
    - `name`: is the cronjob name.
    - `url`: is the url callback to called.
    - `expression`: is cron expression format.
    - `status`: tell if the cronjob is active or paused.
    - `retries`: the number of attempts to send cronjob.
    - `request_timeout`: is the retry timeout.
    - `stop_at`: is the expire time of the cronjob.

#### `POST` `/cronjobs`
Create a new cronjob.
- Method: `POST`
- Endpoint: `/cronjobs`
- Input:
    The `Content-Type` HTTP header should be set to `application/json`

    ```json
   {
      "name":"cronjob-1",
      "url":"your-api/job",
      "expression": "0 5 * * * *",
      "status": 1,
      "retries": 2,
      "request_timeout": 3,
      "stop_at": "2019-11-11 01:00:00",
   }
    ```
- Responses:
    * 201 Created
    ```json
   {
      "id": 1,
      "name":"cronjob-1",
      "url":"your-api/job",
      "expression": "0 5 * * * *",
      "status": 1,
      "retries": 2,
      "request_timeout": 3,
      "stop_at": "2019-11-11 01:00:00",
      "updated_at": "2016-12-10 14:02:37",
      "created_at": "2016-12-10 14:02:37"
   }
    ```
    * 422 Unprocessable entity:
    ```json
    {
      "status":"invalid_cronjob",
      "message":"<reason>"
    }
    ```
    * 400 Bad Request
    ```json
    {
      "status":"invalid_json",
      "message":"Cannot decode the given JSON payload"
    }
    ```
    Common reasons:
    - the cronjob job already scheduled. The `message` will be `Cronjob already exists`
    - the expression must be crontab format.
    - the retry must be between `0` and `10`, max 10 times
    - the status must be `1` or `9`
    - request_timeout: callback request timeout, default 3s
    - url: callbackurl, will fill a param "__cronId=$id" in the url, for distinguish which job cause. 
    - stop_at: as expire time, job will stop after that time.

#### `GET` `/cronjobs/:id`
Get a specific cronjob.
- Method: `GET`
- Endpoint: `/cronjobs/:id`
- Responses:
    * 200 OK
    ```json
   {
      "name":"cronjob-1",
      "url":"your-api/job",
      "expression": "0 5 * * * *",
      "status": 1,
      "retries": 2,
      "request_timeout": 3,
      "stop_at": "2019-11-11 01:00:00",
      "updated_at": "2016-12-10 14:02:37",
      "created_at": "2016-12-10 14:02:37"
   }
    ```
    * 404 Not Found
    ```json
    {
      "status":"cronjob_not_found",
      "message":"The cronjob was not found"
    }
    ```

#### `DELETE` `/cronjobs/:id`
Remove a scheduled cronjob.
- Method: `DELETE`
- Endpoint: `/cronjobs/:id`
- Responses:
    * 200 OK
    ```json
    {
      "status":"cronjob_deleted",
      "message":"The cronjob was successfully deleted"
    }
    ```
    * 404 Not Found
    ```json
    {
      "status":"cronjob_not_found",
      "message":"The cronjob was not found"
    }
    ```

#### `PATCH` `/cronjobs/:id`
Update a cronjob.
- Method: `PATCH`
- Endpoint: `/cronjobs/:id`
- Input:
    The `Content-Type` HTTP header should be set to `application/json`

    ```json
   {
      "expression": "0 2 * * * *"
   }
    ```
- Responses:
    * 200 OK
    ```json
   {
      "name":"cronjob-1",
      "url":"your-api/job",
      "expression": "0 2 * * * *",
      "status": 1,
      "retries": 2,
      "request_timeout": 3,
      "stop_at": "2019-11-11 01:00:00",
      "updated_at": "2016-12-10 14:02:37",
      "created_at": "2016-12-10 14:02:37"
   }
    ```
    * 404 Not Found
    ```json
    {
      "status":"cronjob_not_found",
      "message":"The cronjob was not found"
    }
    ```
    * 422 Unprocessable entity:
    ```json
    {
      "status":"invalid_json",
      "message":"Cannot decode the given JSON payload"
    }
    ```
    * 400 Bad Request
    ```json
    {
      "status":"invalid_cronjob",
      "message":"<reason>"
    }
    ```

## Cron Format
The cron expression format allowed is:

|Field name| Mandatory?|Allowed values|Allowed special characters|
|:--|:--|:--|:--|
|Seconds      | Yes        | 0-59            | * / , -|
|Minutes      | Yes        | 0-59            | * / , -|
|Hours        | Yes        | 0-23            | * / , -|
|Day of month | Yes        | 1-31            | * / , - ?|
|Month        | Yes        | 1-12 or JAN-DEC | * / , -|
|Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?|
more details about expression format [here](https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format)


## response promise by callbacked server

* message format:
    * json is expected.
    * code value cause the job status changing:
    * case code equal 200, job will continue;
    * case code equal 400, job will stop, do not run callback in next cron cycle;
    * message will be ignored.
    ```json
    {
        "code":200,
        "message":"OK"
    }
    ```
  
## TODO LIST
- [x] cycle log, counter each job.
- [x] replace http handler by gin
- [ ] api authentication
- [ ] to be distributed
- [ ] web ui console

