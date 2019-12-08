# Rest-Server-With-Postgres

Simple Rest API Server with Postgres Database

## Prerequisites

Install GO dependencies

```
. ./resource/installDependency.sh
```

Install Postgres, Create User & Database & Table ,and Insert some Test Data

| id | name | price |
| --- | --- | --- |
| 1 | apple | 10 |
| 2 | banana | 20 |
| 3 | orangeorange | 30 |

```
. ./resource/createDatabase.sh
```

## Running the program

```
go run main.go
```

## API

### HTTP Get Request to Query item with {id}
```
. "/items/{id}"
```

### HTTP Post Request to Insert item with {name} & {price}
```
. "/items/{name}/{price}"
```

### HTTP Delete Request to Query item with {id}
```
. "/items/{id}"
```
