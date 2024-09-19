# Scootin Aboot 

An app to simulate a scooter rental company.

## Getting started

To set the application up you can simply use this two command after successfully cloning the repository:

```aqua
    go mod download
    docker compose up
```

You would need to have docker installed on your local machine. Running the second command command
will run the application on port 8081 and setup basic data in Redis. It will also
set up and simulate three mobile clients that are randomly using, riding and freeing
the scooters for the first seconds of the application life. After this time this predefined
clients will stop using the app anymore.

The app will be still accessible on port 8081, so to play with it further use the OpenAPI spec enabled in the <b>/docs</b> folder.
If you don't use a known clientID the app would return with 403 Forbidden status so add header:

```aqua
Client-Id: cd81ed3b-c1a5-43f5-b524-35eaebf0430c
```

calling it with curl, postman, etc.

The logs of the application are printed to stdout, so they can be viewed inside the <b>app</b> container.
You can view the logs using:

```aqua
    docker logs <app-container-id>
```

where <app-container-id> can be obtained using:

```aqua
    docker ps
```

## Tests

Unit tests are tagged with <i>unit</i> tag. To run the tests use:
```aqua
go test --tags unit ./...
```

## Example of usage
Once you deployed the application in docker containers using instructions from the above paragraph you should be able to connect to the 
application listening on the port 8081. To do so we can use curl. To get all scooters in Ottawa in a range in a shape of rectangle with the
middle of rectangle defined by it's longitude and latitude and the width and the height of the rectangle defined in meters use:
```aqua
curl -X GET \
-H "Client-Id: cd81ed3b-c1a5-43f5-b524-35eaebf0430c" \
"http://localhost:8081/api/v1/scooters?city=Ottawa&longitude=73.55&latitude=45.5&height=20000.0&width=25000.0"
```

You can also add optional availability query param to the request above to filter the scooters by their current status by adding to the end of the above url:
```aqua
&availability=true
```
or
```aqua
&availability=false
```

If you then want to rent one of the scooters obtained in the result pick one of them that has availability set to true and use:
```aqua
curl -X POST \
-H "Client-Id: cd81ed3b-c1a5-43f5-b524-35eaebf0430c" \
-H "Content-Type: application/json" \
-d '{"UUID": "{scooter_uuid}", "longitude": {scooter_longitude}, "latitude": {scooter_latitude}, "city":"Ottawa"}' \
http://localhost:8081/api/v1/rent
```
scooter_uuid is an UUID to identify scooter, scooter_longitude is its longitude and scooter_latitude is its latitude obtained in the GET 
call above.

If you check the logs of the docker container running the application you can notice that the scooter was rented and its localisation is 
tracked.

Now if you want to free the scooter you have just rented use:

```aqua
curl -X POST \
-H "Client-Id: cd81ed3b-c1a5-43f5-b524-35eaebf0430c" \
-H "Content-Type: application/json" \
-d '{"UUID": "{scooter_uuid}"}' \
http://localhost:8081/api/v1/free
```

If you check the logs of the docker container now you should notice that the scooter was successfully freed and the tracking process has 
ended.
