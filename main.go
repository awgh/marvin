package main

import (
	"flag"

	"github.com/awgh/marvin/marvin"
)

func main() {

	var configDir, chainFile, mcflyFile, mackerFile string
	flag.StringVar(&configDir, "confdir", "config", "Config Directory")
	flag.StringVar(&chainFile, "chain", "markov.chain", "Markov Chain File")
	flag.StringVar(&mcflyFile, "mcfly", "mcfly.chain", "McFly Chain File")
	flag.StringVar(&mackerFile, "macker", "macker.chain", "macker Chain File")
	flag.Parse()

	marvin.Run(configDir, chainFile, mcflyFile, mackerFile)
}
