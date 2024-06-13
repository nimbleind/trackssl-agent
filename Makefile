

agent: *.go
	@go build -o agent *.go
	@GOOS=windows GOARCH=amd64 go build -o agent.exe *.go

run: agent
	@./agent

clean:
	@rm -f agent agent.exe

