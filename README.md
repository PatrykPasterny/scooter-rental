# Scootin Aboot


## Getting started

To set the application up you can simply use:


```aqua
    docker compose up
```

You would need to have docker installed on your local machine. Running this command
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

## Architecture


I have chosen the modular monolith architecture as it is easier to build from scratch, easier to maintain
at the beginning of the process of creating app's prototype and app itself. If needed it is easy to break
apart the modules and create microservices out of it enabling better scaling and adding more complexity.


For the modules and resources I have chosen:
- Redis for database, as it has amazing functionalities for geospatial indexing, is fast to read and easy to maintain.
  highly available, consistent and really easy to scale. The usage of NoSQL database here also enables
  easy way to change the data models without necessity of migrations.
  Redis connections are hidden behind ScooterRepository interface, so in case of future decisions regarding database vendor we can
  easily swap it with different implementation of the interface without a need of change in other places of application (apart
  from main.go of course where we set up the application)
- Rental Service uses ScooterRepository to Rent and Free the scooters. It also asks Tracker service to track the scooters
  on their journeys.
- Tracker Service triggered by the Rental Service it runs and stops the process of tracking scooters. In the production env
  this kind of service may be integrated with the scooters' soft itself, so we can use the GPS transmitter of scooter for the updates
  and change the Rental Service call that triggers Tracker Service into a Message Queue event (RabbitMQ or Kafka can be used for it).
  It would loosen the binding between two services.

Assumptions:
- For the UpdateAvailability redis call I decided that we don't want to allow the availability of scooter to be changed
  from false to false and from true to true. To make it happen I used Watch redis command. That also makes the process of renting
  scooters more reliable and easier, because two users can not change the availability to 0 (rent the scooter) at the same time.

Tradeoffs:
- Redis is not a SQL database, so using indexes or more advanced ways of querying data is not available, so we rely on key-value store.
