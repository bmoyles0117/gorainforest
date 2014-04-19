## Overview
This is a golang client library for [RainforestQA](https://www.rainforestqa.com/). RainforestQA automates the tedious QA process with the precision of manual testing!

## Installation
To install gorainforest, simply run `go get github.com/bmoyles0117/gorainforest`.

## Run Tests Example

	package main

	import (
		"github.com/bmoyles0117/gorainforest"
	)

	func main() {
		clientToken := "ABC"
		client := rainforest.NewRainforest(clientToken)

		client.RunTests([]int{1,2})
	}