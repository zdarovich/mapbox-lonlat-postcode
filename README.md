## How to run tests?
```shell
go test ./...
```
## Decisions
- If some goroutine catches an error, parent goroutine decides how to handle an error.
- If there is an error, goroutine pipeline stops reading from stdin and waits till all workers finish

## Not completed
- It is possible to add more test coverage, even TDD tests

## The following command should be used:

```
 cat coordinates.txt | ./your-program -token "api token" -workers "pool size flag" > output.txt
```

