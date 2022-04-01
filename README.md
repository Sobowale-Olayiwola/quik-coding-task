# quik-coding-task

## Task Description
The goal is to write a JSON API in Golang to get the balance and manage credit/debit  operations on the players wallets. For example, you might receive calls on your API to get  the balance of the wallet with id 123, or to credit the wallet with id 456 by 10.00 €. The  storage mechanism to use will be MySQL. 

## Endpoints
The following endpoints are exposed 
Baseurl = localhost
Port = 8080
### Fetches the wallet balance of a particular registered player
* GET 
    * /api/v1/wallets/{wallet_id}/balance 
### Credits the wallet of a particular registered player on a given wallet id
* POST 
    * /api/v1/wallets/{wallet_id}/credit 
### Debits the wallet of a particular registered player on a given wallet id
* POST 
    * /api/v1/wallets/{wallet_id}/debit 
### Registers a player to Quik.
This endpoint is required in order to implement the authorization middleware and associate a specific player to a wallet

* POST 
    * /api/v1/players
### Fetches a player registered on quik
* GET 
    * /api/v1/players/{id}

### Updates a player registered on quik
* PUT 
    * /api/v1/players/{id}

### Deletes a player registered player on quik
* DELETE 
    * /api/v1/players/{id}

### Login a registered player on quik
* POST 
    * /api/v1/players/login
The API documentation can be visited on postman to interact with endpoints to display the JSON response and sample error codes
https://www.postman.com/bold-desert-829444/workspace/quik

## To run and application
* cd /cmd/api
* go build
Or use air a development tool to generate an executable file in tmp dir

## To run tests
This test coverage is focused on the business logic for crediting and debiting players wallet

* cd /wallet/service
* go test -v
### coverage 
* cd /wallet/service
* go test cover
This result into a coverage of 76.3% covering all the edge cases of the business logic

