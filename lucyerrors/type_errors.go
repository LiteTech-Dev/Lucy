package lucyerrors

import "errors"

var NoLucyError = errors.New("lucy is not installed, run `lucy init` before downloading mods")

var InvalidPlatformError = errors.New("invalid platform")
var PackageSyntaxError = errors.New("invalid package syntax")
var EmptyPackageSyntaxError = errors.New("empty package string")
