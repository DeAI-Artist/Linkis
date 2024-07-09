# MintAI: A Blockchain L1-L2 Ecosystem Unlocking the Potential of Decentralized AI

![banner](docs/tendermint-core-image.jpg)

[![Version][version-badge]][version-url]
[![API Reference][api-badge]][api-url]
[![Go version][go-badge]][go-url]
[![License][license-badge]][license-url]


Tendermint Core includes both Layer-1 (L1) and Layer-2 (L2) implementations. L1 coordinates and records network activities, securely replicating them across many machines. Meanwhile, L2 contracts handle the storage and transportation of value, including coins and tokens in all kinds of types and formats.

For protocol details, refer to the [MintAI WhitePaper](https://mintai.gitbook.io/whitepaper/).

For a detailed analysis of the consensus framework in the L1 layer, including BFT safety and liveness proofs, refer to the paper, "[The latest gossip on BFT consensus](https://arxiv.org/abs/1807.04938)" and its corresponding [Tendermint GitHub repo](https://github.com/tendermint/tendermint).


_NOTE: This is only the dev version of MintAI core, both in the L1 and L2 implementations, for testnet purposes. As the project progresses, we are excited to include more features in both layers and anticipate major upgrades compared to the current versions. We warmly welcome any kind of contributions! For more information, see [our contribution policy](SECURITY.md)._

## Quick Start
### Overview
The network is composed of validators, service providers (miners), and clients. To earn network rewards, you need to either participate as a service provider or a validator.

### Validators
Validators are responsible for securing the network by proposing and validating new blocks. To set up a validator node, follow these steps:

#### Install MintAI:
```shell
git clone https://github.com/DeAI-Artist/MintAI.git
cd MintAI
make build
make install
```
#### Initialize the validator node:
```shell
mintai init
```
#### Start the validator node:
```shell
mintai node
```
From here your validator node is up and running :)

> Note: The current implementations listen to all incoming transactions through the p2p module. To run a local node for testing purposes, modify the configuration file located at `.mintai/config/.config.toml`. Inside the file, change `laddr = "tcp://0.0.0.0:26656"` to `tcp://127.0.0.1"` to restrict network access to localhost only.

#### Send transactions:
The general transaction format in the MintAI network is:
```shell
curl -s 'RPC_laddr:26657/broadcast_tx_commit?tx="TX_CONTENT"'
```
For example, if we register a client with the `client_name` field set to `Client1`, have one of the validators' `RPC_laddr` as `134.209.85.94`, and set `TX_CONTENT` as `0000003b7b22636f6e74...`, then:
```shell
curl -s '134.209.85.94:26657/broadcast_tx_commit?tx="0000003b7b22636f6e74656e74223a2265794a6a62476c6c626e5266626d46745a534936496b4e73615756756444456966513d3d222c2274797065223a317dffef7045be04d7de12b68871a6837c67e26a0d89ca4af124bf0b449041f12a6b25f4d388962f0f7e704cfc93f8a47ba0dc9b15118819cf2067ecd9542b17b4251c"'
```
with the expected output as follows:
```text
{"jsonrpc":"2.0","id":-1,"result":{"check_tx":{"code":0,"data":null,"log":"","info":"","gas_wanted":"0","gas_used":"0","events":[],"codespace":"","sender":"","priority":"0","mempoolError":""},"deliver_tx":{"code":0,"data":null,"log":"","info":"","gas_wanted":"0","gas_used":"0","events":[{"type":"app","attributes":[{"key":"dHJhbnNhY3Rpb24gc2VuZGVy","value":"MHg2YzI1YjcyQ0Q2ODA3RDEwNjc4QjQ1N0I2RTYzRkI3OTNhZTAzMEVi","index":true}]}],"codespace":""},"hash":"39EB7DB89963695E683D862C48C1EB167C757E74E435D6B0C3505A12FE4B4686","height":"88083"}}
```
If `check_tx` and `deliver_tx` have a `code` of `0`, then the transaction is successfully validated and processed, with its content being included in the global states.
We can validate this by using the `query` API on another validator in the network with the command below:
```shell
curl -s '178.128.168.223:26657/abci_query?data="clientRegistration_0x6c25b72CD6807D10678B457B6E63FB793ae030Eb"'
```
To send transactions that can be passed into the mempool, one needs to formulate the transaction according to the messagetype and transaction format specified in this documentation.
#### Transaction fees:
Unlike interacting with mainnets like Solana or Ethereum, MintAI L1 transactions are completely **free of gas fees**, allowing clients to navigate the network without any economic friction.
One can refer to the [MintAI protocol](https://arxiv.org/pdf/2310.19099) to understand how this is accomplished in a secure and robust manner.
## Minimum requirements

| Requirement | Notes             |
|-------------|-------------------|
| Go version  | Go 1.18 or higher |


## Contributing

Please abide by the [Code of Conduct](CODE_OF_CONDUCT.md) in all interactions.

Before contributing to the project, please take a look at the [contributing
guidelines](CONTRIBUTING.md) and the [style guide](STYLE_GUIDE.md). You may also
find it helpful to read the [specifications](./spec/README.md), and familiarize
yourself with our [Architectural Decision Records
(ADRs)](./docs/architecture/README.md) and
[Request For Comments (RFCs)](./docs/rfc/README.md).

## Versioning

### Semantic Versioning

Tendermint uses [Semantic Versioning](http://semver.org/) to determine when and
how the version changes. According to SemVer, anything in the public API can
change at any time before version 1.0.0

To provide some stability to users of 0.X.X versions of Tendermint, the MINOR
version is used to signal breaking changes across Tendermint's API. This API
includes all publicly exposed types, functions, and methods in non-internal Go
packages as well as the types and methods accessible via the Tendermint RPC
interface.

Breaking changes to these public APIs will be documented in the CHANGELOG.

### Upgrades

In an effort to avoid accumulating technical debt prior to 1.0.0, we do not
guarantee that breaking changes (ie. bumps in the MINOR version) will work with
existing Tendermint blockchains. In these cases you will have to start a new
blockchain, or write something custom to get the old data into the new chain.
However, any bump in the PATCH version should be compatible with existing
blockchain histories.

For more information on upgrading, see [UPGRADING.md](./UPGRADING.md).

### Supported Versions

Because we are a small core team, we only ship patch updates, including security
updates, to the most recent minor release and the second-most recent minor
release. Consequently, we strongly recommend keeping Tendermint up-to-date.
Upgrading instructions can be found in [UPGRADING.md](./UPGRADING.md).

## Resources

### Libraries

- [Cosmos SDK](http://github.com/cosmos/cosmos-sdk); A framework for building
  applications in Golang
- [Tendermint in Rust](https://github.com/informalsystems/tendermint-rs)
- [ABCI Tower](https://github.com/penumbra-zone/tower-abci)

### Applications

- [Cosmos Hub](https://hub.cosmos.network/)
- [Terra](https://www.terra.money/)
- [Celestia](https://celestia.org/)
- [Anoma](https://anoma.network/)
- [Vocdoni](https://docs.vocdoni.io/)

### Research

- [The latest gossip on BFT consensus](https://arxiv.org/abs/1807.04938)
- [Master's Thesis on Tendermint](https://atrium.lib.uoguelph.ca/xmlui/handle/10214/9769)
- [Original Whitepaper: "Tendermint: Consensus Without Mining"](https://tendermint.com/static/docs/tendermint.pdf)
- [Tendermint Core Blog](https://medium.com/tendermint/tagged/tendermint-core)
- [Cosmos Blog](https://blog.cosmos.network/tendermint/home)

## Join us!

Tendermint Core is maintained by [Interchain GmbH](https://interchain.berlin).
If you'd like to work full-time on Tendermint Core,
[we're hiring](https://interchain-gmbh.breezy.hr/)!

Funding for Tendermint Core development comes primarily from the
[Interchain Foundation](https://interchain.io), a Swiss non-profit. The
Tendermint trademark is owned by [Tendermint Inc.](https://tendermint.com), the
for-profit entity that also maintains [tendermint.com](https://tendermint.com).

[bft]: https://en.wikipedia.org/wiki/Byzantine_fault_tolerance
[smr]: https://en.wikipedia.org/wiki/State_machine_replication
[Blockchain]: https://en.wikipedia.org/wiki/Blockchain
[version-badge]: https://img.shields.io/github/tag/tendermint/tendermint.svg
[version-url]: https://github.com/DeAI-Artist/MintAI/releases/latest
[api-badge]: https://img.shields.io/badge/API-Online-brightgreen
[api-url]: https://pkg.go.dev/github.com/DeAI-Artist/MintAI
[go-badge]: https://img.shields.io/badge/go-1.21-blue.svg
[go-url]: https://github.com/moovweb/gvm
[discord-badge]: https://img.shields.io/discord/669268347736686612.svg
[discord-url]: https://discord.gg/cosmosnetwork
[license-badge]: https://img.shields.io/badge/License-GPL--3.0-lightgreen
[license-url]: https://github.com/DeAI-Artist/MintAI/blob/main/LICENSE
[sg-badge]: https://sourcegraph.com/github.com/DeAI-Artist/MintAI/-/badge.svg
[sg-url]: https://sourcegraph.com/github.com/DeAI-Artist/MintAI?badge
[tests-url]: https://github.com/DeAI-Artist/MintAI/actions/workflows/tests.yml
[tests-badge]: https://github.com/DeAI-Artist/MintAI/actions/workflows/tests.yml/badge.svg?branch=main
[lint-badge]: https://github.com/DeAI-Artist/MintAI/actions/workflows/lint.yml/badge.svg
[lint-url]: https://github.com/DeAI-Artist/MintAI/actions/workflows/lint.yml
