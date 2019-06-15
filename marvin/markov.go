package marvin

import (
	"math/rand"
	"time"

	"github.com/awgh/markov"
)

var markovChains []*markov.Chain

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator for markov.
}

// InitMarkovChains - Loads an array of markov chains from given array of filenames
func InitMarkovChains(chains ...string) {
	// Markov Chain setup
	markovChains = make([]*markov.Chain, len(chains))
	for i, v := range chains {
		markovChains[i], _ = loadChainFile(v)
	}
}
