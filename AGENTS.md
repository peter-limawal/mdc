# MDC

## Go

Run before committing Go changes:

```bash
gofmt -w <changed-go-files>
go vet ./...
go test ./...
```

## Naming

- `command`: command arguments (`[]string`)
- `commandJSON`: serialized command data
- `execCmd`: `*exec.Cmd` subprocess
- `job`: input job
- `gotJob`: job returned by an operation
- `storedJob`: job retrieved to verify persistence

## Tests

- Test observable behavior with the standard `testing` package.
- Keep tests explicit and behavior-focused.
- Add only the smallest test that covers a new non-trivial behavior.
