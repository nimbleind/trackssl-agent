
agent: *.go
	@GOOS=darwin GOARCH=amd64 go build -o trackssl-agent-mac *.go
	@GOOS=windows GOARCH=amd64 go build -o trackssl-agent-windows.exe *.go
	@GOOS=linux GOARCH=amd64 go build -o trackssl-agent-linux *.go

test:
	@go test -v ./...

run: agent
	@./agent

clean:
	@rm -f trackssl-agent-mac \
		trackssl-agent-windows.exe \
		trackssl-agent-linux \
		trackssl-agent.tar.gz \
		trackssl-agent.zip

release: agent
	git archive --format=tar.gz --output=trackssl-agent.tar.gz HEAD
	git archive --format=zip --output=trackssl-agent.zip HEAD

