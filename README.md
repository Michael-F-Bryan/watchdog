# Watchdog Service

Checks if the CMT website and Confluence are still up.

## Get Started

Gets a local version of the Watchdog Service.
```bash
$ git clone https://github.com/Curtin-Motorsport-Team/watchdog \
	$GOPATH/src/github.com/Curtin-Motorsport-Team/watchdog
```

Change into the newly created directory.

Complie the Watchdog Service.
```bash
$ go build
```

Run the Watchdog Service.
Make the `\` or `/` is correct for your OS.
```bash
.\watchdog
```

## Flag

### timeout

Can be used to change the amount of time before the service check will timeout.
This time in in seconds with the default being 10 seconds.

Usage:
```bash
$ .\watchdog.exe -timeout=10
```
or:
```bash
go run main.go -timeout=10
```
