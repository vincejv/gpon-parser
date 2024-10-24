# GPON Parser (GPON Stats exporter)

Supports the following ONT models

* FiberHome HG6245D (Globe Telecom Philippines firmware)
  * `ONT_MODEL: hg6245d_globe`
* FiberHome AN5506_04F1A (Globe Telecom Philippines firmware) and other generic FH ONT with generic firmware
  * `ONT_MODEL: an5506_stock`
* ZTE F670
  * `ONT_MODEL: zte_f670`
* ZTE F660
  * `ONT_MODEL: zte_f670`
* ZLT G3000A (Globe Telecom Philippines firmware)
  * `ONT_MODEL: zlt_g3000a`

## Environment variables
* `ONT_WEB_HOST`
  * IP address of ONT
  * Default: depends on modem
* `ONT_WEB_PORT`
  * Port on which ONT Web UI is listening to
  * Default: depends on modem
* `ONT_WEB_PROTOCOL`
  * Web protocol which the ONT web gui uses, typically set as `http` or `https`
  * Default: `http`
* `ONT_WEB_USER`
  * ONT Web UI username
  * Default: depends on modem
* `ONT_WEB_PASS`
  * ONT Web UI password
  * Default: depends on modem
* `ONT_TELNET_PORT`
  * ONT Web UI password
  * Default: `23`
* `ONT_POLL_SEC`
  * Specifies the frequency on how often the GPON stats are pulled from the ONT
  * Default: `60`
* `LISTEN_PORT`
  * Port on which the exporter listens to
  * Default: `8092`
* `LISTEN_IP`
  * Ip address on which the exporter listens to
  * Default: `0.0.0.0`

## Running
Docker Pull
```sh
docker pull vincejv/gpon-parser:latest
```
Docker Run
```sh
docker run -d \
  --name gpon-parser \
  --restart unless-stopped \
  vincejv/gpon-parser:latest
```
Docker Compose
```yaml
version: '3'

services:
  gpon-parser:
    image: vincejv/gpon-parser:latest
    container_name: gpon-parser
    restart: unless-stopped
    environment:
      ONT_MODEL: "zte_f670"
```

## REST API Paths
`/gpon/allInfo`
```json
{
  "deviceStats": {
    "memoryUsage": 54.885117384596136,
    "cpuUsage": 1.31,
    "cpuDtlUsage": [
      0.1,
      2.52
    ],
    "deviceModel": "F660",
    "modelSerial": "FHTTXXXXXX",
    "softwareVersion": "V1.1.20P3N6B",
    "uptime": 86673
  },
  "opticalStats": {
    "rxPower": -26.5757,
    "txPower": 2.7781,
    "temperature": 44,
    "supplyVoltage": 3.229,
    "biasCurrent": 13.5
  }
}
```
`/gpon/deviceInfo`
```json
{
  "memoryUsage": 54.880947416704885,
  "cpuUsage": 2.4749999999999996,
  "cpuDtlUsage": [
    0.1,
    4.85
  ],
  "deviceModel": "F660",
  "modelSerial": "FHTTXXXXXX",
  "softwareVersion": "V1.1.20P3N6B",
  "uptime": 86748
}
```
`/gpon/opticalInfo`
```json
{
  "rxPower": -26.5757,
  "txPower": 2.7781,
  "temperature": 44,
  "supplyVoltage": 3.229,
  "biasCurrent": 13.55
}
```

## Footnotes

### Building the package
```sh
go build -ldflags "-s -w"
```

### Running
```sh
go run .
```

### ARM Build on Windows
```powershell
$env:GOARCH='arm'
$env:GOOS='linux'
```

### ARM Build on Linux
```sh
export GOARCH='arm'
export GOOS='linux'
```