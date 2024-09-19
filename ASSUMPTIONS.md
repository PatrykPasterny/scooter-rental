# Assumptions

## Architecture

I have chosen the modular monolith architecture as it is easier to build from scratch, easier to maintain
at the beginning of the process of creating app's prototype and app itself. If needed it is easy to break
apart the modules and create microservices out of it enabling better scaling and adding more complexity.


For the modules and resources I have chosen:
- Redis for database, as it has amazing functionalities for geospatial indexing, is fast to read and easy to maintain,
  highly available, consistent and really easy to scale. The usage of NoSQL database here also enables
  easy way to change the data models without necessity of migrations. It can also be very well scaled localisation wise having for
  example one redis instance per city/region.
- Redis connections are hidden behind ScooterRepository interface, so in case of future decisions regarding database vendor we can
  easily swap it with different implementation of the interface without a need of change in other places of application (apart
  from main.go of course where we set up the application)
- Rental Service uses ScooterRepository to Rent and Free the scooters. It also asks Tracker service to track the scooters
  on their journeys.
- Tracker Service triggered by the Rental Service runs and stops the process of tracking scooters using ScooterRepository. In the production
  env this kind of service may be integrated with the scooters' soft itself, so we can use the GPS transmitter of scooter for the updates
  and change the Rental Service call that triggers Tracker Service into a Message Queue event (RabbitMQ or Kafka can be used for it).
  It would loosen the binding between two services.
- Fake clients are run as separate docker container.

## Other
Other than strictly architecture assumption from my side:
- For the UpdateAvailability redis call I decided that we don't want to allow the availability of scooter to be changed
  from false to false and from true to true. To make it happen i needed to peek on the latest availability value and update it when the
  condition is valid. I used Watch redis command with MULTI to be sure it has transaction like behaviour. That also makes the process of
  renting scooters more reliable and easier, because two users can not change the availability to false (rent the scooter) at the same time.
- Scooters does not communicate with the API, instead the tracker service does which is written in a way it could be transferred to
  scooters software as mentioned in the architecture part and use scooters GPS device to update location, so the scooter in this approach
  does not have to authenticate.

##Tradeoffs
The main tradeoff assigned with the current approach are:
- Redis is not a SQL database, so using indexes or more advanced ways of querying data is not available, so we rely on key-value store.
