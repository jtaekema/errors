# Errors

A simple way to add tracing information to an error.

[![Go Report Card](https://goreportcard.com/badge/github.com/jtaekema/errors)](https://goreportcard.com/report/github.com/jtaekema/errors) [![Build Status](https://travis-ci.org/jtaekema/errors.svg?branch=master)](https://travis-ci.org/jtaekema/errors) [![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](https://github.com/jtaekema/errors/blob/master/LICENSE)

## How to use

To add tracing to an error, you would do something like:

    return errors.Trace(err)

To add a trace with some context use, do this:

    return errors.Tracef(err, "trying to perform action...")

The context is displayed after the file/line number and can be very useful for debugging by just glancing at the error trace.

## Example

    package main

    import (
      "fmt"

      "github.com/jtaekema/errors"
    )

    func internal() error {
      return errors.New("this is an error")
    }

    func doSomethingComplex() error {
      if 1 != 5 {
        err := internal()
        if err != nil {
          return errors.Tracef(err, "internal didn't work when '1 != 5'")
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

The error details would look like:

    /Users/jtaekema/go/src/github.com/jtaekema/errors/app/app.go:10 [error] this is an error
    /Users/jtaekema/go/src/github.com/jtaekema/errors/app/app.go:17 [error] internal didn't work when '1 != 5'

## Best Practices

This is somewhat up to you, but I think there are two philosophies.

1. Always add a trace when you return an error.
2. Only add a trace when you can add context.
