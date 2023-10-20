Download necessary go modules using the `go.mod` file

```
go mod download

```

1. Generate appropriate abis from the `abis` folders

```
make compile
```

2. RPC urls are given as a multiline env variable, refer to `.env.example` for reference

3. Edit and run `main.go` as required
