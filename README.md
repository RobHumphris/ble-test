# ble-test
A set of tests for the Go ble library `github.com/go-ble/ble`

## Making a snap
The following steps were followed to make this into a snap.
1. `mkdir snap`
2. `mkdir ble-test`
3. `cd ble-test`
4. `snapcraft init`
5. `vim snap/snapcraft.yaml`

```yaml
name: ble-test
version: git
summary: Bluetooth test
description: |
  Bluetooth test

grade: devel
confinement: devmode

parts:
  ble-test:
    source: .
    plugin: go
    go-importpath: github.com/RobHumphris/ble-test

apps:
  ble-test:
    command: ble-test
```

6. `git clone https://github.com/RobHumphris/ble-test.git`
7. `cd ble-test`