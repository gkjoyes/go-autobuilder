# go-autobuilder
Build go projects automatically when files get modified.

## Install

```sh
    go get -u github.com/george-kj/go-autobuilder
```

## Features

- Build automatically when files get modified.

- Set environment variables for our project by just passing the `env` file path

- Run custom commands while building the project. Example, pass `-race` flag for race condition detection.

- Run custom commands while `running`. Example, this can support CLI interfaces tools like a [cobra](https://github.com/spf13/cobra).

- Run custom commands before each build, Example, run Go formatting tool [gofmt](https://golang.org/cmd/gofmt/) before each build.

- We can pass multiple commands with each option by separating with single space.

- If a single command has any spaces need to escape this using a backward slash. eg:- `gofmt\ -w\ .`

## Usage

- Update `PATH` with `GOPATH/bin`

    ```sh
    Usage: go-autobuilder
    options:
        -p 	The directory to be watch.
        -n	Project name.
        -e	Environment file path.
        -v	Prints the version.
        -b	Build only mode.
        -cc	Custom commands to run before the build.
        -bc	Custom commands to run while building.
        -rc	Custom commands to run while running.
    ```