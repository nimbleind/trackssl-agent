# TrackSSL Agent

## MacOS Development

Execution:

```
TRACKSSL_AUTH_TOKEN=*authtoken* TRACKSSL_AGENT_TOKEN=*agenttoken* ./agent
```

## Windows Development

Building:

```
go build -o agent.exe
```

Execution:

```
cmd /V /C "set TRACKSSL_AUTH_TOKEN=*authtoken*&&set TRACKSSL_AGENT_TOKEN=*agenttoken*&&agent.exe"
```


