# Backendify Solution

The Backendify Solution is my idea of how the problem should be solved in a way that I would typically address service design.

Not only in a functional way to solve the problem, but with a design pattern I believe to be best. 

## Description

This project is defined as a Webservice with 2 endpoints. The first of which is a simple health check returning a 200 Status code everytime
in order verify the service is running. The Second endpoint provides the main business logic of the service allowing the caller to request company 
data after providing a company ID and country ISO code. The project also communicates with other backend services to complete 
the request. These backend address correspond to various ISO codes and are unique to each country code. The addresses are provided as command line
arguments when running the service. 

### Architectural Design

The design architecture follows the Service Layer Pattern to create modularity between the various functionality.
For instance the Transport Layer focuses solely on the communication to and from the caller whether that be a RESTful or 
RPC design, this work is abstracted away from the lower layer. Further down in this implementation we have the Service Layer
which is where the data validation and business logic occurs. This layer then communicates with the various regional backends
which in turn act as the Storage layers. 


![Architecture Design.png](assets%2FArchitecture%20Design.png)


One of the key advantages of this design pattern is the interdependence of the different layers. This allows for the ability 
to make larger technological changes to different aspects such as the Database or communication protocol without requiring 
code changes to other layers.

## Technical Specs

- `/status`: A simple health check returning a 200 Status code everytime.
  - Method: GET
- `/company`: Provides the company resource upon request.
  - Method: GET
  - Query Requirements: The endpoint expects an `id`(string) and `country_iso`(2 letter string).
    - format: `/company?id=XX&country_iso=YY`
  - Response Object: JSON
```json
{
"id": "string, the company id requested by a customer",
"name": "string, the company name, as returned by a backend",
"active": "boolean, indicating if the company is still active according to the active_until date",
"active_until": "RFC 3339 UTC date-time expressed as a string, optional."
}
```

## Usage
To use the service you may run the following commands or the Makefile commands in the next section.

### Build
Since Go is a compiled langauge you must first compile the code into executable with the following 'build' command:

`go build -o .`

### Run
The program requires the address of the various backend address with their corresponding ISO codes to be passed
in as any number of arguments to the run command. The arguments should be in the format '{ISO_Code}={Backend_Address}'

Example: `ru=http://localhost:9001`

When you are ready to run simply call the executable followed by the provided arguments.
`./backendify {ARGS}`

Example: `./backendify uk=http://localhost:9001 ru=http://localhost:9002`

### Test
You may run the automated unit tests with the following command:

`go test`

## Makefile
### Build
`make build`

### Run
`make run`

### Testing
`make test`