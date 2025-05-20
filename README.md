# Errors

A simple way to add tracing information and context to an error.

[![Go Report Card](https://goreportcard.com/badge/github.com/jtaekema/errors)](https://goreportcard.com/report/github.com/jtaekema/errors) [![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](https://github.com/jtaekema/errors/blob/master/LICENSE)

## How to use

This is intended to be a drop in replacement for the `errors` package included in the standard library with additional functionality.

To create a new error with tracing information:

    errors.New("something bad happened")

or for your convience

    errors.New("action: %v failed to complete successfully", action)

To add tracing and context to an error, you would do something like:

    return errors.Wrap(err, "trying to perform some action...")

or

    return errors.Wrapf(err, "trying to perfom: %v", action)

The context is displayed after the file/line number and can be very useful for debugging by just glancing at the error trace.

## Example

    package main

    import (
      "fmt"

      "github.com/jtaekema/errors/v2"
    )

    func internal() error {
      return errors.New("this is an error")
    }

    func doSomethingComplex() error {
      if 1 != 5 {
        err := internal()
        if err != nil {
          return errors.Wrap(err, "internal didn't work when '1 != 5'")
        }
      }
      return nil
    }

    func main() {
      err := doSomethingComplex()
      if err != nil {
        fmt.Println(errors.Details(err))
      }
    }

The printed error details would look like:

    /app/main.go:10 [error] this is an error
    /app/main.go:17 [error] internal didn't work when '1 != 5'

## Best Practices

This is somewhat up to the reader, but I believe it is best to always add context when you return an error.
