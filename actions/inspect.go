package actions

import (
	"errors"

	"github.com/zajicekn/pokedexcli/pokedex"
)

func InspectPokemon(name string) (*pokedex.Pokedex, error) {
	// Create an instance of the pokedex
	dex := pokedex.GetGlobalDex()

	// Check to see if pokemon name is in the user's pokedex
	_, ok := dex.Data[name]
	if !ok {
		return nil, errors.New("you have not caught that pokemon")
	}

	return dex, nil
}
