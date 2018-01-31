# Usage Store

### CI status
master: [![CircleCI](https://circleci.com/gh/ernestio/usage-store/tree/master.svg?style=shield)](https://circleci.com/gh/ernestio/usage-store/tree/master)  | develop: [![CircleCI](https://circleci.com/gh/ernestio/usage-store/tree/develop.svg?style=shield)](https://circleci.com/gh/ernestio/usage-store/tree/develop)

### Description

Usage store is intended to be a data store to track the lifetime of specific resource types.
The database structure is basically an entry list of different timeframes where the specific entity was up and running.

Even you can interact as usual with this store through [its endpoints](main.go#L14)
```
usage.get
usage.del
usage.set
usage.find
```
the database population is done by listening at events happening on [Ernest](http://ernest.io) nats network. Usage store is listening at [this events](main.go#L48)

#### instance.create.*.done
Creates a new entry for the specific instance starting at the current time

#### instance.delete.*.done
Closes an existing entry

#### instance.update.*.done
Depending on the field `powered` [true|false] it will execute the same logic on the previous two, in case powered is not set it will skip any process.


## Installation

```
make deps
make install
```

## Running Tests

This repository has integration tests described [here](https://github.com/ernestio/ernest/blob/master/internal/features/cli/usage_report.feature), they will be compiled on the repo linked CI. If you want to run them locally, follow instructions described on [ernest repo](https://github.com/ernestio/ernest)

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 ernest.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
