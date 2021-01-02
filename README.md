# StagelinQ for Golang

This library implements Denon's StagelinQ protocol, allowing any Go application to talk to devices that are compatible with this protocol on the network.

## Stability

The code of this project is an **experimental** reverse-engineering effort and therefor can behave erratically in untested cases. Currently, this code only has been practically tested with the Denon Prime 4.

If you have any other Denon devices you would like to test this library against, please do! Even better, you can let me know if you run into any bugs by reporting them [as an issue ticket](https://github.com/icedream/go-stagelinq/issues).

## Testing

This project uses Go tests, they can be run with this command:

    go test ./...

## License

This code is licensed under the MIT license. For more information, please read [LICENSE](LICENSE).
