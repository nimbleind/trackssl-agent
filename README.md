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

### Decrypting windows password.

```
aws lightsail get-instance-access-details --instance-name Windows_Server_2022-1 | \
    jq -r '.accessDetails.passwordData.ciphertext' | \
    base64 -d | \
    openssl pkeyutl -decrypt -inkey ~/development/nimble/Nimble_Metabase.pem; echo
```


