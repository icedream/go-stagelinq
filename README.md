# StagelinQ for Golang

[![Go Reference](https://pkg.go.dev/badge/github.com/icedream/go-stagelinq.svg)](https://pkg.go.dev/github.com/icedream/go-stagelinq)

This library implements Denon's StagelinQ protocol, allowing any Go application to talk to devices that are compatible with this protocol on the network.

An example application is provided that, if running successfully, will output information like this:

![Screenshot of the example CLI](docs/screenshot.png)

## Features

- Automatically discover StagelinQ-compatible devices on the network
- Access state map information such as currently playing track metadata, fader values, etc.

## Stability

The code of this project is an **experimental** reverse-engineering effort and therefor can behave erratically in untested cases. Currently, this code only has been practically tested with the Denon Prime 4.

If you have any other Denon devices you would like to test this library against, please do! Even better, you can let me know if you run into any bugs by reporting them [as an issue ticket](https://github.com/icedream/go-stagelinq/issues).

## Testing

This project uses Go tests, they can be run with this command:

    go test ./...

## License

This code is licensed under the MIT license. For more information, please read [LICENSE](LICENSE).
