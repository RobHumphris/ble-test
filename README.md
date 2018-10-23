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
8. `snapcraft`
9. `sudo snap install ble-test_0+git.cb54c62_amd64.snap --devmode --dangerous`
10. Registered the snap name ble-snap `https://dashboard.snapcraft.io/register-snap/`
11. `snapcraft login`
12. `snapcraft push ble-test_0+git.cb54c62_amd64.snap`

To get the snap running on the Dell Edge Gateway 3001 (obviously the version numbers and ip addresses will change...)
1. `scp ./ble-test_0+git.cb54c62_amd64.snap admin@192.168.1.95:~`
2. Login to the Gateway `ssh admin@192.168.1.95`
3. `sudo snap install ble-test_0+git.cb54c62_amd64.snap --devmode --dangerous`
4. `sudo ble-test`