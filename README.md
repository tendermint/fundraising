# Fundraising Module

The fundraising module is a Cosmos SDK module that that serves fundraising feature, which provides an oppotunity for new projects to onboard the ecosystem. It does not only allow projects to raise funds, but also increase their brand awareness before launching their projects. 

The fundraising module is built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

- see the [main](https://github.com/tendermint/fundraising/tree/main) branch for the latest 
- see [releases](https://github.com/tendermint/fundraising/releases) for the latest release

## Dependencies

If you haven't already, install Golang by following the [official docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

Requirement | Notes
----------- | -----------------
Go version  | Go1.16 or higher
Cosmos SDK  | v0.44.0

## Installation

```bash
# Use git to clone the source code and install `fundraisingd`
git clone https://github.com/tendermint/fundraising.git
cd fundraising
make install
```

## Getting Started

To get started to the project, visit the [TECHNICAL-SETUP.md](./TECHNICAL-SETUP.md) docs.

## Documentation

The fundraising module documentation is available in [docs](./docs) folder and technical specification is available in [specs](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/README.md) folder. 

These are some of the documents that help you to quickly get you on board with the fundraising module.

## Contributing

We welcome contributions from everyone. The [main](https://github.com/tendermint/fundraising/tree/main) branch contains the development version of the code. You can branch of from main and create a pull request, or maintain your own fork and submit a cross-repository pull request. If you're not sure where to start check out [CONTRIBUTING.md](./CONTRIBUTING.md) for our guidelines & policies for how we develop fundraising module. Thank you to all those who have contributed to fundraising module!
