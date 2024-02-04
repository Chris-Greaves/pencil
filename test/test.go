// Package of shared functions to be used in Unit Tests
package test

import (
	"runtime"
)

func FileNotFoundText() string {
	if runtime.GOOS == "windows" {
		return "The system cannot find the file specified"
	} else {
		return "no such file or directory"
	}
}

func PathNotFoundText() string {
	if runtime.GOOS == "windows" {
		return "The system cannot find the path specified"
	} else {
		return "no such file or directory"
	}
}
