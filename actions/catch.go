package actions

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/zajicekn/pokedexcli/pokestructs"
)

func GetPokemon(area *pokestructs.Config) (*pokestructs.ShallowPokemon, error) {
	// This will be where we check the pokedex to see if we already have
	// the pokemon we are trying to catch

	// Determine the URL being used to get the batch of pokemon
	var url string
	if area.Next != "" {
		url = area.Next
	} else {
		url = "https://pokeapi.co/api/v2/pokemon"
	}

	// Make a GET request to receive response
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error Rquest: ", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error Reading Response: ", err)
	}

	// Unmarshal the data into JSON
	var pokemonJSON *pokestructs.ShallowPokemon
	err = json.Unmarshal(responseData, &pokemonJSON)
	if err != nil {
		log.Fatal("Error Unmarshalling: ", err)
	}

	// Update the config struct to process next pokemon batch
	area.Next = pokemonJSON.Next
	area.Previous = pokemonJSON.Previous

	return pokemonJSON, nil
}

func CatchPokemon(area *pokestructs.Config, name string) (*pokestructs.PokemonStats, error) {
	// Initialize an empty url string to soon hold the URL
	// with appropriate pokemon ID attached
	url := ""

	// Continuously loop over the locations of the shallowJSON.
	// Find the URL that matches the name given as a parameter
	flag := true
	for flag {
		// Call Map to get the locations
		shallowPokemon, err := GetPokemon(area)
		if err != nil {
			log.Fatal(err)
		}

		// Another for loop to iterate over the names of the areas
		// If it matches the name of the parameter, break out of loop
		for _, result := range shallowPokemon.Results {
			if result.Name == name {
				url = result.URL
				flag = false
				break
			}
		}
		if !flag {
			area.Next, area.Previous = "", nil
			break
		}

		// Update config with new pagination URLs
		area.Next = shallowPokemon.Next
		area.Previous = shallowPokemon.Previous

		if area.Next == "" {
			return nil, errors.New("that pokemon doesn't exist")
		}
	}

	// Make a GET request for the individual pokemon
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Request Error: ", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error Reading: ", err)
	}

	// Write the JSON data into the new pokemonJSON variable
	var pokemonJSON *pokestructs.PokemonStats
	err = json.Unmarshal(responseData, &pokemonJSON)
	if err != nil {
		log.Fatal("Error Unmarshal: ", err)
	}

	return pokemonJSON, nil
}
