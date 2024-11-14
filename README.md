# Supervise

Dirt simple (close to stupid) web GUI for [supervisord](http://supervisord.org/).

## Motivation

Needed a web GUI for supervisord, and was not clever enough to realize it already exist.

## Building

```sh
go build .
```

### Optional

```sh
sudo mv supervise /usr/bin/
```

## Configuration

You need to have a configuration file in order to run the application.

On Linux and MacOS, the default config file is located here: `/etc/supervise/settings.json`

For Windoze you have to use the command-line flag `-config` to specify where the config file is located. The reason is that I have no idea where a good default place for configuration files would be on a Windoze system.

You need to have a configuration file in order to specify the usernames and passwords of the users given access:

```json
{
	"accounts": {
		"donald": "duck"
	}
}
```

In this example, the user "donald" has the password "duck"

## Running the application

```sh
./supervise
```

Or, if you installed the binary to your PATH:

```sh
supervise
```

### Command line options

* `-addr` - _string_ - Serve on this address. Default: `127.0.0.1:9988`
* `-config` - _string_ - Path to config file. Default: `/etc/supervise/settings.json` (on Linux and MacOS)
* `-dev` - _bool_ - Dev mode. Will not use embedded files. Default: `0`

The `-dev` flag is meant to be used during development. The `supervise` binary contains the static files needed for the web server. This is called "embedding", making it easy to deploy anywhere, but during development you would need to rebuild the whole application for every single change to the (for example) CSS. Easier to just pass `-dev=0` instead.

## Running supervise using supervisord

Example configuration:

```ini
[program:supervise]
command=supervise
stdout_logfile=/var/log/supervisord/supervise.log
stderr_logfile=/var/log/supervisord/supervise.error.log
```

This example assumes that you have installed `supervise` on your PATH.

If you choose to not install `supervise` on your PATH (and I don't blame you, not at all), then something like this would be proper:

```ini
[program:supervise]
command=/path/to/supervise
stdout_logfile=/var/log/supervisord/supervise.log
stderr_logfile=/var/log/supervisord/supervise.error.log
```

### Running as root

Depending on your system, you'll need to ensure that you run `supervise` as root:

```ini
[program:supervise]
command=/path/to/supervise
user=root
stdout_logfile=/var/log/supervisord/supervise.log
stderr_logfile=/var/log/supervisord/supervise.error.log
```

## Third-party libraries

Thank you [gin](https://github.com/gin-gonic/gin), and thank you [jQuery](https://jquery.com/) ❤
