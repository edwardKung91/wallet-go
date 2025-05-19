# This Code was developed in windows environment

# Decisions
- Included a user_id as a part of the wallet table to show that we eventually want to include an owner for each wallet
- All the ids, including waller and user ids, are UUIDs because usually banks/fintechs enforce a format for their account numbers/ids. I chose to use UUIDs as there are generators available online
- The APIs should ideally also ideally include ways to provide Ids from upstream instead of just creating its own and streaming it out but as a Phase 1 take home assignment i chose to go with simple implementations
- Idempotency keys should ideally have been created as well to prevent duplicate requests and allow for retry functionality but as an initial submission I chose to make the transactions simple
- Creating DB level locks or optimistic locking mechanism should are also something I chose to skip for now for simplicity

# Reviewers
```
wallet-go
| - cmd
| - |
| - | - server
| - | - |
| - | - |- main.go -> "This handles service initialisation"
| - |
| - pkg -> "All service related files and components are here"
| - |
| - | - config
| - | - |
| - | - | - config.go -> "This loads all the DB related configurations from a .env file"
| - |
| - | - db
| - | - |
| - | - | - schema
| - | - | - |
| - | - | - | - schema.sql -> "This contains the table creation SQL queries for use in the service"
| - | - |
| - | - | - db_init.go -> "This contains functions for initialising the db and the db tables"
| - | - |
| - | - | - postgres.go -> "This initialises the DB and also returns the DB connection to be used by the service"
| - |
| - | - router
| - | - |
| - | - | - router.go -> "This contanins the service and handler instanciation and the definition of external APIs"
| - |
| - | - wallet -> "Contains all files relating to the service itself
| - | - |
| - | - | - constants.go -> "Contains definition of constants used within the service"
| - | - |
| - | - | - errors.go -> "contains definition of errors used with in the service"
| - | - |
| - | - | - handler.go -> "contains the logic to process each request and pass it on to an appropriate backend function"
| - | - |
| - | - | - handler_test.go -> "test for the main handlers"
| - | - |
| - | - | - mock_service.go -> "contains mock services for the handler tests"
| - | - |
| - | - | - model.go -> "contains struct definitions for structures used by the service
| - | - |
| - | - | - service.go -> "contains the backend logic that needs to be performed to service each request"
| - | - |
| - | - | - service_test.go -> "tests for the main service functions"
| 
| - .env -> "contains values for DB configuration"
|
| - ASSUMPTIONS.md -> "contains assumptions and also file directory for reviewers
|
| - IMPROVEMENTS.md -> "contains potential improvements to the existing code"
|
| - README.md -> "contains instructions on how to use the APIs and how to run the server"
```

# Time spent
- 12hrs total