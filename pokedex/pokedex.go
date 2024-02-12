package pokedex

import (
	"sync"

	"github.com/zajicekn/pokedexcli/pokestructs"
)

type Pokedex struct {
	Data map[string]pokestructs.PokemonStats
}

var (
	globalDex *Pokedex
	once      sync.Once
)

// Create a globalDex instance for functions to use
func InitGlobalDex() {
	once.Do(func() {
		// Initialize
		globalDex = NewDex()
	})
}

// GetGlobalDex returns the singleton instance of the pokedex
func GetGlobalDex() *Pokedex {
	return globalDex
}

// This function creates a new dex for the request data
func NewDex() *Pokedex {
	newDex := Pokedex{
		Data: make(map[string]pokestructs.PokemonStats),
	}
	return &newDex
}
