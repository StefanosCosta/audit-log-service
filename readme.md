# Audit log service

This is an audit log service as per the requirements of Canonical's interview.

## Storage Decision:

As event types are open ended out data is unstructured. This initially points to a nosql solution that will be able to take any kind of json payload and store it as part of the database. The first thought was to use mongodb which would be perfect for this. Nevertheless, this requires a hefty setup process and a requirement of the interview is to:

Provide instructions for deploying the solution using a single command on Ubuntu;

I explored other Nosql solutions with mongodb that would require less setup, as well as ones that would be best for audit log services:

https://aws.amazon.com/compare/the-difference-between-redis-and-mongodb/#:~:text=Redis%20is%20an%20in%2Dmemory,provides%20speed%20and%20data%20durability.
https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
https://aws.amazon.com/nosql/in-memory/
https://stackoverflow.com/questions/27543134/redis-handling-of-huge-dataset
https://yuriktech.com/2021/07/19/Auditing-in-Microservices/
https://aws.amazon.com/redis/#:~:text=Redis%20offers%20a%20fast%2C%20in,desktop%20users%20at%20a%20time.
https://yuriktech.com/2021/07/19/Auditing-in-Microservices/

Despite this, none of the solutions seemed appropriate as all of them required some form of lenghty setup.
Finally, I decided to emulate such solution using sqlite3, as it would be an in memory database in the file system and would not require any installation of docker or special packages. A schema with the common fields of the event and an extra json field is used to make this possible in db/eventsRepository/events.go. The service is write-intensive, and as sqlite's website boasts handling 100000 requests a day, 20% of which write to the database I assumed that this would be adequate. With a max file size of 281 TB, running out of storage shouldn't be an issue either. The speed trade-off is that marshalling and unmarshalling JSON from sqlite3 is likely much slower than a NoSQL database.

## Service design:

The service has to accept event data sent by other systems and provide an HTTP endpoint for querying recorded event data by field values. To accept data by other systems, I chose to implement a Remote Procedure Call(RPC) which is usually used in distributed systems. However, I soon discovered that the only way to call an rpc server using curl, which is required by the document, is to use json-rpc. Therefore two servers are deployed: the rpc server uses json-rpc to receive requests listening at port 5001, and an http server listens at port 8000 for registering, login, and event querying.

Authentication is made using rsa keys to generate and verify jwt authentication tokens that are passed to the user upon login, and then need to be used in subsequent requests' authorisation header to be able to query or create new events. There is also a register service for registering a new user. Even though an IAM authentication service would be the proper choice here, the given the constraints of this project, make this impossible.

The yaml and key file would also normally be stored on Amazon S3 or equivalent. The login service in this file is also redundant for rpc calls, as usually in microservices there is an authorization service dedicated to doing so, and the rpc calls are being made by a service that is already authorized to do so.

The functionality is split into two packages:
authservice and eventservice.
As most of the functionality exists there, unit tests are made to cover that code. Unit tests include gomock and database cleaning.

The event service is currently set to accept requests that do not contain the common fields just in case there is an event that is like that, as event structure received is not pre-defined.

## CURL Instructions

Sample requests for each api are available in the payloads.txt file in this repository. It is simply a matter of adding the appropriate parameter and jwt token.

## How to run:

A simple "go run main.go" should do the trick.

The repository can be found at https://github.com/StefanosCosta/audit-log-service
