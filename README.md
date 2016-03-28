[![Build Status][circleci-image]][circleci-url]

# godbpf

The **godbpf** project is a library to manipulate Simcity 4 DBPF files written in Go.

These files package a collection of files into a single archive, which can be read and loaded into the game at run time.

## Installation

More  to come...

## Usage

More to come...

## Testing

The entire library can be tested from its root directory using the go test command
`go test`

### Code Coverage

To measure code coverage from testing, the go test command must be run from each
package
`go test -coverprofile=coverage.out`

To view the coverage reports in HTML form, use the go tool cover command
`go tool cover -html=coverage.out`

## Contributing

If you are interested in contributing to this project, please the [CONTRIBUTING.md](CONTRIBUTING.md) document
for guidelines and practices used on this project.

## License

This project is made available under the MIT license.  For complete license terms, please refer to the [LICENSE.md](LICENSE.md) file.

[circleci-image]: https://circleci.com/gh/marcboudreau/godbpf.svg?style=shield&circle-token=6275c851b1ca8b2191032fcda36ebe6dcdf8f640
[circleci-url]: https://circleci.com/gh/marcboudreau/godbpf
