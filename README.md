# test task for Enlabs

To check:
1. Prepare the database (PostgresQL).
2. Run cmd/createenv/createenv.go to prepare the test environment.
   We do not have user management, so we are creating a special user for the integration test.
   Test user configured in types/const.go.

run:   
`cd cmd/createenv`   
`go run createenv.go -c "postgres://docker:docker@127.0.0.1/test1?sslmode=disable"`

3. Run unit and integration tests:   
`cd ../..`   
`go test ./...`  



