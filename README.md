# Watchdog Service

Checks if the CMT website and Confluence are still up.

## Get Started

To get a local version of the Watchdog Service.
```bash
$ git clone https://github.com/Curtin-Motorsport-Team/watchdog \
	$GOPATH/src/github.com/Curtin-Motorsport-Team/watchdog
```

Change into the newly created directory as used above.

Compile the Watchdog Service.
```bash
$ go build
```

Now you can run the Watchdog Service,
Be aware that you make sure the `\\` or `/` is correct for your OS.

```bash
.\watchdog
```

## Flags

### timeout

Can be used to change the amount of time before the service check, will timeout.
This amount of time is in seconds, with the default being 10 seconds.

This effects all service checks.

Usage:
```bash
$ .\watchdog.exe -timeout=10
```
or:
```bash
go run main.go -timeout=10
```


## Testing

At the moment the test suite is still largely incomplete, however you can run
the tests using the following command:

```bash
$ go test
```
