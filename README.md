![Cover](./assets/cover.png)
`go-herder` is a tiny programs runner. This utility runs your programs and scripts and provides an API for managing them and tracking states.

## Installation


## Documentation

### Beginning
To run `go-herder`, two configuration files are required: `go-herder.yml` and `go-herder.db` (sqlite).
These are the default names that you can change using the `-yml` and `-db` flags when launching the utility.
File `.yml` contains go-herder configurations, and a file `.db` contains the parameters of the programs being run.

To create these files, use the `go-herder init` command.

### Configure

### Running

### API
#### /herder
    /run
    /state
    /kill
#### /herder/processes/:id
    /run
    /state
    /kill