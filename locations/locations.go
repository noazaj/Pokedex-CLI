package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zajicekn/pokedexcli/pokecache"
	"github.com/zajicekn/pokedexcli/pokestructs"
)

func GetLocationMap(area *pokestructs.Config) (*pokestructs.ShallowLocation, error) {
	// Create cache instance
	cache := pokecache.GetGlobalCache()

	// Determine the URL to use for the request
	var url string
	if area.Next != "" {
		url = area.Next
	} else {
		url = "https://pokeapi.co/api/v2/location-area"
	}

	// Check to see if the URL is in the cache
	cachedData, found := cache.Get(url)
	if found {
		// Unmarshal the cached data
		var locationData *pokestructs.ShallowLocation
		err := json.Unmarshal(cachedData, &locationData)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		// Use the unmarshalled data and return early
		for _, result := range locationData.Results {
			fmt.Printf("Name: %s\n", result.Name)
		}
		return locationData, nil
	}

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Request Error: ", err)
	}
	defer resp.Body.Close()

	// Read and unmarshal the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Reading Error: ", err)
	}

	var locationJSON *pokestructs.ShallowLocation
	err = json.Unmarshal(responseData, &locationJSON)
	if err != nil {
		log.Fatal("Error Unmarshal: ", err)
	}

	// Update config with new pagination URLs
	area.Next = locationJSON.Next
	area.Previous = locationJSON.Previous

	// Marshal the ShallowLocation and add it to the cache
	marshalledData, err := json.Marshal(locationJSON)
	if err != nil {
		fmt.Println("Error Marshal: ", err)
	}
	cache.Add(url, marshalledData)

	return locationJSON, nil
}

func GetLocationMapb(area *pokestructs.Config) (*pokestructs.ShallowLocation, error) {
	// Use cache instance
	cache := pokecache.GetGlobalCache()

	// Determine the URL to use for the request
	var url string
	if *area.Previous != "" {
		url = *area.Previous
	} else {
		return &pokestructs.ShallowLocation{}, errors.New("error: no previous url")
	}

	// Check to see if the URL is in the cache
	cachedData, found := cache.Get(url)
	if found {
		// Unmarshal the cached data
		var locationData *pokestructs.ShallowLocation
		err := json.Unmarshal(cachedData, &locationData)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		// Use the unmarshalled data
		for _, result := range locationData.Results {
			fmt.Printf("Name: %s\n", result.Name)
		}
		return locationData, nil
	}

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Request Error: ", err)
	}
	defer resp.Body.Close()

	// Read and unmarshal the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Reading Error: ", err)
	}

	var locationJSON *pokestructs.ShallowLocation
	err = json.Unmarshal(responseData, &locationJSON)
	if err != nil {
		log.Fatal("Error Unmarshal: ", err)
	}

	// Update config with new pagination URLs
	area.Next = locationJSON.Next
	area.Previous = locationJSON.Previous

	// Marshal the ShallowLocation and add it to the cache
	marshalledData, err := json.Marshal(locationJSON)
	if err != nil {
		fmt.Println("Error Marshal: ", err)
	}
	cache.Add(url, marshalledData)

	return locationJSON, nil
}

func ExploreLocation(area *pokestructs.Config, name string) (*pokestructs.Location, error) {
	// Create a cache instance
	cache := pokecache.GetGlobalCache()

	// Initialize an empty url string to soon hold the URL
	// with appropriate ID attached
	url := ""

	// Start the URL from the beginning
	area.Next = "https://pokeapi.co/api/v2/location-area"

	// Continuously loop over the locations of the shallowJSON.
	// Find the URL that matches the name given as a parameter
	flag := true
	for flag {
		// Call Map to get the locations
		shallowJSON, err := GetLocationMap(area)
		if err != nil {
			log.Fatal(err)
		}

		// Another for loop to iterate over the names of the areas
		// If it matches the name of the parameter, break out of loop
		for _, result := range shallowJSON.Results {
			if result.Name == name {
				url = result.URL
				flag = false
				break
			}
		}
		if !flag {
			break
		}

		// Update config with new pagination URLs
		area.Next = shallowJSON.Next
		area.Previous = shallowJSON.Previous
	}

	// Check to see if name with url is in cache
	cachedData, found := cache.Get(url)
	if found {
		// Unmarshal the cached data
		var pokemonJSON *pokestructs.Location
		err := json.Unmarshal(cachedData, &pokemonJSON)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		// Print out marshalled data and return early
		for _, val := range pokemonJSON.PokemonEncounters {
			fmt.Printf(" - %s\n", val.Pokemon.Name)
		}
		return pokemonJSON, nil
	}

	// If not in cache, make a normal GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Request Error: ", err)
	}
	defer resp.Body.Close()

	// Read the resp and unmarshal
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error Reading: ", err)
	}

	// Write the JSON data into the new pokemonJSON variable
	var pokemonJSON *pokestructs.Location
	err = json.Unmarshal(responseData, &pokemonJSON)
	if err != nil {
		log.Fatal("Error Unmarshal: ", err)
	}

	// Marshal the pokemon data and add it to the cache
	marshalPokemonData, err := json.Marshal(pokemonJSON)
	if err != nil {
		log.Fatal("Error Marshalling: ", err)
	}
	cache.Add(url, marshalPokemonData)

	return pokemonJSON, nil
}
