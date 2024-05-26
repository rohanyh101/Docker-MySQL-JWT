# MySQL JWT Docker CRUD Microservice
This repository contains a microservices implementation with three services: users, tasks, and projects.
Each service has its endpoints and all endpoints are protected with JWT authentication. 
The architecture uses MySQL for database management and Docker for containerization, ensuring secure, scalable, and easily deployable services.

## Run Instructions
- Ensure that Docker is installed on your system.
- Run `make docker-build` to build your Docker image with the defined credentials.
- Run `make docker-run` to start your MySQL server.
- Run `make run` to start the project on http://localhost:3000.
- (Optional) To run the tests, execute `make test`.

finally don't forget to test all the endpoints in `Postman` or `ThunderClient`

Adios... 👋



