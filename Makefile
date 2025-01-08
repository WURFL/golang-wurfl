test:
	go test ./... -run ^Test -v -timeout 30m -race -shuffle=on

benchmarks:
	go test ./... -bench Benchmark_ -run=^$ go
