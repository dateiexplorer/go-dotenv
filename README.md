# go-dotenv

A dotenv package for the Go programming language.

This module provides functions to load variables form a file into the environment.
It is heavily influcenced by the original Ruby implementation
(https://github.com/bkeepers/dotenv) and various other Go libraries, e.g.
https://github.com/joho/godotenv.

Other than the most existing libraries this module provides a well defined
interface giving the ability to implement further interpreters (so called
'Shells') to support individual syntax.
A term 'Shell' was chosen because the overall goal of this interface was to
have the ability to support various shell syntax, e.g. bash, zsh, fish,
PowerShell, etc.

However, the interface is heavily extensible so that you can implement a
'Shell' with your own syntax or to support other file formats, e.g. TOML or
JSON.

The current package implements only a Bash shell (making use of the 'echo'
command of a real Bash) and a so called Basic shell, which implements a subset
of the Bash functionality. This implementation is also cross platform.

## Installation

To use this module you can simply download it in your Go environment with the
following command.

```
go get github.com/dateiexplorer/go-dotenv
```

## Usage

After you installed the module you can simply add e.g. a configuration file
`.env` in the root of your project with the following lines:

```
SECRET_KEY=YOURSECRETKEY
I_AM_COOL=$USER is cool!
```

Then in your Go application you can do something like this:

```go
package main

import (
    "log"
    "os"

    "github.com/dateiexplorer/go-dotenv"
)

func main() {
    err := dotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    secretKey := os.Getenv("SECRET_KEY")
    msg := os.Getenv("I_AM_COOL")

    // Now you can whatever you suppose to do with that vairables
    log.Println(msg)
    log.Printf("Never log your secret key: %v\n", secretKey)
}

```

You can also specify the name of your env file as well as the Shell you want
to use for loading variables, e.g.

```go
_ = godotenv.LoadWith(shell.Bash, "my_custom.env")
```

## Contributing

Oh yeah! Contribute, give it to me... wait, what?
Contributions are very welcome of course.

Please create a PullRequest with your new features. Only well tested
contributions will be accepted!

## Special thanks

So, also that this code has another code base than joho's version, this project
is heavily influenced by the ideas from https://github.com/joho/godotenv.

Maybe this library is also capable for your use case.