
We use https://github.com/mattes/migrate for migrations.

To install use 
```
go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli
```
