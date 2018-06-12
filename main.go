package main

import (
	"github.com/Wazzymandias/blockstack-crawler/cmd"
	"runtime"
)

func main() {
	// recommended value for improved DB performance on SSDs,
	// see: https://groups.google.com/forum/#!topic/golang-nuts/jPb_h3TvlKE/discussion
	runtime.GOMAXPROCS(128)

	cmd.Execute()
}
