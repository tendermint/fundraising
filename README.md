# Fundraising Module

The fundraising module is a Cosmos SDK module that serves the fundraising feature that provides an oppotunity for new projects to onboard into the Cosmos ecosystem. The fundraising module allow projects to raise funds and increase their brand awareness before launching their projects. 

The fundraising module is built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

- [main](https://github.com/tendermint/fundraising/tree/main) branch for the latest 
- [releases](https://github.com/tendermint/fundraising/releases) for the latest release

## Dependencies

If you haven't already, install Golang by following the official Go [install docs](https://golang.org/doc/install). Make sure that your `GOPATH` and `GOBIN` environment variables are properly set up.

Requirement | Notes
----------- | -----------------
Go          | version 1.16 or higher
Cosmos SDK  | v0.44.0

## Installation

```bash
# Use git to clone the source code and install `fundraisingd`
git clone https://github.com/tendermint/fundraising.git
cd fundraising
make install

# Install binary in testing mode enables MsgAddAllowedBidder to add an allowed bidder 
make install-testing
```

## Getting started

To get started with the project, visit the [TECHNICAL-SETUP.md](./TECHNICAL-SETUP.md) docs.

## Documentation

The fundraising module documentation is available in [docs](./docs) folder and technical specification is available in [specs](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/README.md) folder. 

These documents can help you start using the fundraising module.

## Contributing

We welcome contributions from everyone. The [main](https://github.com/tendermint/fundraising/tree/main) branch contains the development version of the code. You can branch of from main and create a pull request, or maintain your own fork and submit a cross-repository pull request. If you're not sure where to start check out [CONTRIBUTING.md](./CONTRIBUTING.md) for our guidelines and policies for how we develop the fundraising module. Thank you to everyone who has contributed to the fundraising module!
