# Qume Sample Client for Authenticated WebSocket Connection

This is a sample WebSocket client application, written in Go, to exemplify how a client should interact with any WebSocket endpoint provided by Qume that requires authentication.

## Tool Configuration

It is necessary to pass an *YAML* file with configuration values to the tool in order to execute it. The path to such file should be provided to the tool via `CONFIG_FILE` enironment variable.

```yaml
connection:
  host: <Qume WebSocket Host>
  endpoint: <Qume WebSocket Endpoint>
  origin: <WebSocket Client Origin>
authentication:
  key: <API Key for User>
  secret: <API Secret for User>
  password: <User Password>
```

## Building and Running
```bash
git clone git@github.com:qume-exchange/ws-private-example.git
cd wsclient
docker build -t wsclient .
# Make sure to fullfil config/config.yml with correct data
docker run -v $PWD/config:/config -e CONFIG_FILE=/config/config.yml wsclient
```
