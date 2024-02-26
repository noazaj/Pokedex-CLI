package pokeAPI

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/zajicekn/Pokedex-CLI/actions"
	"github.com/zajicekn/Pokedex-CLI/locations"
	"github.com/zajicekn/Pokedex-CLI/pokedex"
	"github.com/zajicekn/Pokedex-CLI/pokestructs"
)

func Map(area *pokestructs.Config) error {
	// Call GetLocationMap to make request
	locationJSON, err := locations.GetLocationMap(area)
	if err != nil {
		log.Fatal(err)
	}

	// Update config with new pagination URLs
	area.Next = locationJSON.Next
	area.Previous = locationJSON.Previous

	// Print out the location-areas
	for _, result := range locationJSON.Results {
		fmt.Printf("Name: %s\n", result.Name)
	}
	return nil
}

func Mapb(area *pokestructs.Config) error {
	locationJSON, err := locations.GetLocationMapb(area)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Update config with new pagination URLs
	area.Next = locationJSON.Next
	area.Previous = locationJSON.Previous

	// Print out the location-areas
	for _, result := range locationJSON.Results {
		fmt.Printf("Name: %s\n", result.Name)
	}
	return nil
}

func Explore(area *pokestructs.Config, name string) error {
	pokemonJSON, err := locations.ExploreLocation(area, name)
	if err != nil {
		log.Fatal(err)
	}

	// Print the Pokemon in the explored area
	fmt.Printf("Exploring %s...\nFound Pokemon:\n", name)
	for _, val := range pokemonJSON.PokemonEncounters {
		fmt.Printf(" - %s\n", val.Pokemon.Name)
	}
	return nil
}

func Catch(area *pokestructs.Config, name string) error {
	// Create a Pokedex instance
	dex := pokedex.GetGlobalDex()

	// Need to clear out area struct for pokemon use only now
	area.Next = ""
	area.Previous = nil

	pokemonJSON, err := actions.CatchPokemon(area, name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create a random number to be used as a 'chance' to catch
	// the pokemon based on their BaseExperience
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	if rand.Intn(300) > pokemonJSON.BaseExperience {
		fmt.Printf("%s was caught!\n", name)
		dex.Data[name] = *pokemonJSON
		fmt.Println("You may inspect it with the inspect command.")
		return nil
	}
	fmt.Printf("%s escaped!\n", name)
	return nil
}

func Inspect(name string) error {
	// Create pokemon variable which is of type pokedex
	pokemon, err := actions.InspectPokemon(name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Print all the pokemon information to the user
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %v\n", pokemon.Data[name].Height)
	fmt.Printf("Weight: %v\n", pokemon.Data[name].Weight)
	fmt.Println("Stats:")
	for _, pokemon := range pokemon.Data[name].Stats {
		fmt.Printf("    -%s: %v\n", pokemon.Stat.Name, pokemon.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemon := range pokemon.Data[name].Types {
		fmt.Printf("    -%s\n", pokemon.Type.Name)
	}

	return nil
}

func Pokedex() {
	// Create pokedex instance
	dex := pokedex.GetGlobalDex()

	// Check to ensure there is data in pokedex
	if len(dex.Data) == 0 {
		fmt.Println("No pokemon in your pokedex. Try catching one!")
		return
	}

	// Print the names of the pokemon in the pokedex
	fmt.Println("Your Pokedex:")
	for _, name := range dex.Data {
		fmt.Printf(" - %s\n", name.Name)
	}
}
