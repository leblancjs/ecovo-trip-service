# Trip Service
## Table of Contents
* [Introduction](#introduction)
* [To-Do](#to-do)
* [Configuration](#configuration)
* [Build and Test](#build-and-test)
* [Deploy](#deploy)
* [Endpoints](#endpoints)
* [Errors](#errors)

## Introduction
The trip service implements the trip REST API. It makes it possible to access a trip's details, such as it's departure and arrival information, as well as create and delete a trip.

## To-Do
* Write unit tests for the use case (service) layer
* Add filters for GET (like orderBy, radius, etc.)
* Decide which data we wont to give in a response for driverId and vehicleId

## Configuration
The application's database connection and Auth0 domain are configured using environment variables. To avoid having to define them every time the service is run, they are kept in the `.env` file at the root of the repository.

The table below enumerates the different environment variables.

|Name|Required|Description|
|---|---|---|
|AUTH_DOMAIN|Yes|Domain where the trip info endpoint is hosted (ex. my.domain.com)|
|DB_HOST|Yes|URI to where the database is hosted|
|DB_USERNAME|Yes|Username to use to to establish the database connection|
|DB_PASSWORD|Yes|Password to use to establish the database connection|
|DB_NAME|Yes|Name of the database to use on the server|
|DB_CONNECTION_TIMEOUT|No|Time to wait before giving up on connecting to the database|
|API_KEY|Yes|API key used for google maps API|

## Build and Test
### Prerequisites
#### Docker
Docker is used to simplify the build and test processes. It makes it possible
to build and run the application without needing to install Go, and also makes it much easier to define environment variables to use to configure the service (see the next section).

Please download and install [Docker Desktop](https://www.docker.com/products/docker-desktop), and make sure that it is running on your machine before you proceed.

### Step 1 - Build an Image
In order to run the application locally to test it, we need to build an image
using Docker.

To do so, run following command in a terminal:

```
docker build --tag=trip-service .
```

You will need to rebuild the image every time a change is made in the code, or when new changes are pulled.

Don't worry, it doesn't take that long.

### Step 2 - Run the Image in a Container
To run the service, we need to run the image we built in the previous step in a
container using Docker.

To do so, run the following command in a terminal and replace `<PORT>` with the port you want to use to access the API:

```
docker run -it -p <PORT>:8080 --env-file .env trip-service
```

It is important to note that the `--env-file` argument is used to tell Docker
to define the environment variables found in the `.env` file in the Docker
container. Otherwise, the service will not start.

## Deploy
The service can be deployed to [Heroku](https://heroku.com) by pushing a Docker
image to its container registry, and releasing it in a Heroku application.

### Environment Variables
It is important to note that the service still needs those environment
variables! On Heroku, they need to be defined in the dashboard as Config Vars.
Without them, the service will fail to start.

### Prerequisites
The same prerequisites defined in the Build and Test section apply here.

#### Heroku CLI
The [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli#download-and-install)
is used to deploy the application to Heroku. Please download and install it on your machine.

##### Login
To log in to Heroku, enter the following command in a terminal:

```
heroku login
```

It should open a web browser in which you can log in using the Ecovo account credentials, which can be found on Google Drive.

This step only needs to be done once, after you've installed the Heroku CLI.

##### Login to the Container Registry
In order to push images and release them on Heroku, you need to log in to the
Heroku container registry.

To do so, enter the following command in a terminal:

```
heroku container:login
```

##### Link the Git Repository to the Heroku Application
To make sure that we deploy the service to the right application on Heroku, we
need to link the Git repository to the application.

To do so, run the following command in a terminal:

```
git remote add heroku git@heroku.com:ecovo-trip-service.git
```

This step only needs to be done once, after you've cloned the Git repository.

#### Step 1 - Push the Image to the Container Registry
To build and push the image to the Heroku container registry, use the following command:

```
heroku container:push web
```

#### Step 2 - Release the Container
To release the container that was pushed in the previous step, use the following command:

```
heroku container:release web
```

#### Step 3 - (Optional) Check the Logs
To check the logs to make sure everything went well, use the following command:

```
heroku logs --tail
```

## Endpoints

### GET /trips/{id}
#### URL Parameters
##### id
The trip's unique identifier generated when it is created.

#### Request
##### Headers
```
Authorization: Bearer {access_token}
```

#### Response
##### Status Code
200 OK

##### Headers
```
Content-Type: application/json
```

##### Body
```
{
    "id": {{id}},
    "driverId": {{driverId}},
    "vehicle": {
        "id": {{id}},
        "make": {{make}},
        "year": {{year}},
        "model": {{model}}
    },
    "full": {{full}},
    "leaveAt": {{leaveAt}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
    "arriveBy": {{arriveBy}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
    "seats": {{seats}},
    "stops": [
    	{
    		"id": {{id}},
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	},
            "seats": {{seats}},
            "timestamp": {{timestamp}}
    	},
        {
        	"id": {{id}},
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	},
            "seats": {{seats}},
            "timestamp": {{timestamp}}
        },
        {
            "id": {{id}},
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	},
            "seats": {{seats}},
            "timestamp": {{timestamp}}
    	}
    ],
    "details": {
        "animals": {{animals}},
        "luggages": {{luggages}}
    },
    "reservationsCount": {{reservationCount}},
    "totalTripPrice": {{totalTripPrice}},
	"pricePerSeat": {{pricePerSeat}},
	"totalDistance: {{totalDistance}}
}
```

##### Possible Errors
* 404 Not Found
* 500 Internal Server Error

### GET /trips
#### Query Parameters
##### source (Mandatory)
The trip's source geographic location.

##### destination (Mandatory)
The trip's destination geographic location.

##### radiusThresh
The trip's source or destination threshold in meters.

##### leaveAt
The trip's departure time.

##### arriveBy
The trip's arrival time.

##### seats
The trip's number of seats available.

##### detailsAnimals
The trip's animals allowance.

##### detailsLuggages
The trip's luggages size allowed.

##### driverId
The trip's driver ID (used to get list of trips for a user).

#### Request
##### Headers
```
Content-Type: application/json
Authorization: Bearer {access_token}
```

#### Response
##### Status Code
200 OK

##### Headers
```
Content-Type: application/json
```

##### Body
```
[
    {
        "id": {{id}},
        "driverId": {{driverId}},
        "vehicle": {
            "id": {{id}},
            "make": {{make}},
            "year": {{year}},
            "model": {{model}}
        },
        "full": {{full}},
        "leaveAt": {{leaveAt}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
        "arriveBy": {{arriveBy}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
        "seats": {{seats}},
        "stops": [
            {
                "id": {{id}},
                "point": {
                    "name": {{name}},
                    "longitude": {{longitude}},
                    "latitude": {{latitude}}
                },
                "seats": {{seats}},
                "timestamp": {{timestamp}}
            },
            {
                "id": {{id}},
                "point": {
                    "name": {{name}},
                    "longitude": {{longitude}},
                    "latitude": {{latitude}}
                },
                "seats": {{seats}},
                "timestamp": {{timestamp}}
            },
            {
                "id": {{id}},
                "point": {
                    "name": {{name}},
                    "longitude": {{longitude}},
                    "latitude": {{latitude}}
                },
                "seats": {{seats}},
                "timestamp": {{timestamp}}
            }
        ],
        "details": {
            "animals": {{animals}},
            "luggages": {{luggages}}
        },
        "reservationsCount": {{reservationCount}},
        "totalTripPrice": {{totalTripPrice}},
        "pricePerSeat": {{pricePerSeat}},
        "totalDistance: {{totalDistance}}
    },
]
```

##### Possible Errors
* 404 Not Found
* 500 Internal Server Error

### POST /trips
#### Request
##### Headers
```
Content-Type: application/json
Authorization: Bearer {access_token}
```

##### Body
```
{
    "driverId": {{driverId}},
    "vehicle": {
        "id": {{id}},
        "make": {{make}},
        "year": {{year}},
        "model": {{model}}
    },
    "full": {{full}},
    "leaveAt": {{leaveAt}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
    "arriveBy": {{arriveBy}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
    "seats": {{seats}},
    "stops": [
    	{
    		"point": {
    			"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
    		}
    	},
        {
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	}
        },
        {
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	}
    	}
    ],
    "details": {
        "animals": {{animals}},
        "luggages": {{luggages}}
    }
}
```

#### Response
##### Status Code
* 201 CREATED

##### Headers
```
Content-Type: application/json
```

##### Body
```
{
    "id": {{id}},
    "driverId": {{driverId}},
    "vehicle": {
        "id": {{id}},
        "make": {{make}},
        "year": {{year}},
        "model": {{model}}
    },
    "full": {{full}},
    "leaveAt": {{leaveAt}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
    "arriveBy": {{arriveBy}}, **format : YYYY-MM-DDThh:mm:ss.sZ**
    "seats": {{seats}},
    "stops": [
    	{
    		"id": {{id}},
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	},
            "seats": {{seats}},
            "timestamp": {{timestamp}}
    	},
        {
        	"id": {{id}},
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	},
            "seats": {{seats}},
            "timestamp": {{timestamp}}
        },
        {
            "id": {{id}},
        	"point": {
        		"name": {{name}},
	        	"longitude": {{longitude}},
	        	"latitude": {{latitude}}
        	},
            "seats": {{seats}},
            "timestamp": {{timestamp}}
    	}
    ],
    "details": {
        "animals": {{animals}},
        "luggages": {{luggages}}
    },
    "reservationsCount": {{reservationCount}},
    "totalTripPrice": {{totalTripPrice}},
	"pricePerSeat": {{pricePerSeat}},
	"totalDistance: {{totalDistance}}
}
```

##### Possible Errors
* 400 Bad Request
* 500 Internal Server Error

### DELETE /trips/{id}
#### Required Parameters
##### id (Mandatory)
The trip's unique identifier generated when it is created.

#### Request
##### Headers
```
Authorization: Bearer {access_token}
```

#### Response
##### Status Code
* 200 OK

##### Possible Errors
* 400 Bad Request
* 500 Internal Server Error

## Errors
### Structure
The errors returned by the service have the following format:

```
{
    "code": {code},
    "message": "{message}"
    "requestId": "{requestId}"
}
```

#### Code
The code generally aligns with the HTTP status code. Its purpose is to give a
general idea of what went wrong. As a rule of thumb, if the code is `500`,
something went wrong on the service's end. Otherwise, it's not our fault :D.

#### Message
The message gives additional information related to the error. For example, in
the case of a `400 Bad Request`, it might contain the name of the field that
was missing.

#### Request ID
The request ID is everyone's best friend. When you an error response that has a
`500` status code and an error message that says that you need to contact a
system administrator, you need to keep that ID! If you look at the server logs,
the internal error will be logged with that request ID, so we can find out what
went wrong.

### Possible Errors
|Status Code|Meaning|Description|
|---|---|---|
|400|Bad Request|A bad request could mean that the body is missing a required field, or has an error in its JSON syntax. In the case of a missing field, it should be included in the error message.
|401|Unauthorized|As the name suggests, this means that the user is not authorized to access the resource. Normally, this is because the token is invalid or expired.
|404|Not Found|When no trip can be found for a given ID, we'll tell ya! Try again when it's created ;).
|500|Internal Server Error|We don't like this one. It means that the service made a mistake! It could be that we couldn't encode a response, or that our database flipped us off. Either way, take that precious request ID and ask us to look into it!