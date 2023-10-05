# GPON Parser

Supports the following ONT models

* FiberHome HG6245D (Globe Telecom Philippines firmware)
* FiberHome AN5506_04F1A (Globe Telecom Philippines firmware)
* ZTE F670
* ZTE F660

## Running
Docker Pull
```
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

### Compiling
```
go build -ldflags "-s -w"
```

### ARM Build on Windows
```powershell
$env:GOARCH='arm'
$env:GOOS='linux'
```