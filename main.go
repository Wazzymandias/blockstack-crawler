package main

import (
	"github.com/Wazzymandias/blockstack-profile-crawler/cmd"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(128)

	cmd.Execute()
}
