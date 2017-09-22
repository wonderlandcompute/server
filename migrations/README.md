
Initially install `go get github.com/gemnasium/migrate`.

Create migrations(note `\` before `?`, in zsh at least it is needed):
```
migrate -url postgres://localhost:5432/database_name\?sslmode=disable  create initial
```


To perform migration:
```
migrate -url postgres://localhost:5432/database_name\?sslmode=disable up 1
```