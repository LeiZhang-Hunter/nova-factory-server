package main

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewRunner)
