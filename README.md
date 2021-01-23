# lemonilo-backend 


## Requirement
- install golang
- install redis 
- install postgresql
- install postman


## Documentasi
Step 1 
create database run query : 
```
    CREATE database lemonilo-backend;
```

create table and insert data dummy :
```
 run db in file file\db.sql
```

Step 2
- configuration env 
example : 
env.example

Step 3
```
- Run the server
```
go mod vendor
cd server
go run main.go
```

Step 4
- Configuration Postman
Postman collection : 
    in file lemonilo-backend.postman_collection.json
Postman Env : 
    in file Local.postman_environment.json


## User Dummy 
username = thoriq007 
password = kiasu123