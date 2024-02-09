# spine-go

[![Build Status](https://github.com/enbility/spine-go/actions/workflows/default.yml/badge.svg?branch=dev)](https://github.com/enbility/spine-go/actions/workflows/default.yml/badge.svg?branch=dev)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/enbility/spine-go)
[![Coverage Status](https://coveralls.io/repos/github/enbility/spine-go/badge.svg?branch=dev)](https://coveralls.io/github/enbility/spine-go?branch=dev)
[![Go report](https://goreportcard.com/badge/github.com/enbility/spine-go)](https://goreportcard.com/report/github.com/enbility/spine-go)
[![CodeFactor](https://www.codefactor.io/repository/github/enbility/spine-go/badge)](https://www.codefactor.io/repository/github/enbility/spine-go)

## Introduction

This library provides an implementation of SPINE 1.3 in [go](https://golang.org), which is part of the [EEBUS](https://eebus.org) specification.

Basic understanding of the EEBUS concepts SHIP and SPINE to use this library is required. Please check the corresponding specifications on the [EEBUS downloads website](https://www.eebus.org/media-downloads/).

This repository was started as part of the [eebus-go](https://github.com/enbility/eebus-go) before it was moved into its own repository and this separate go package.

Basic understanding of the EEBUS concepts SHIP and SPINE to use this library is required. Please check the corresponding specifications on the [EEBUS downloads website](https://www.eebus.org/media-downloads/).

__Important:__ In contrast to the EEBUS recommendation to use a "Generic" client feature, this library does not support this for the local device! Instead one should create a feature type with the client role for every required feature.

## Packages

### api

This package contains required interfaces. They are used extensivly to be able to mock everything and implement tests that focus specificaly on a limited set of interface implementations

### integrationtests

This packge contains tests that cover implementations of multiple packages in concert.

### mocks

This package contains auto generated mocks for the interfaces defined in the api package using [Mockery](https://github.com/vektra/mockery).

### model

This package contains the go represenation of the SPINE data model. It makes use of go tags for proper JSON serialization and also for implementing generic SPINE feature to function and data mapping.

### spine

This package contains the implementation for working with the SPINE devices, entites, features, functions and data.

### util

This package contains generic helpers used by most of the packages, e.g. for working with pointers.
