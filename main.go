package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/zajicekn/Pokedex-CLI/pokeAPI"
	"github.com/zajicekn/Pokedex-CLI/pokecache"
	"github.com/zajicekn/Pokedex-CLI/pokedex"
	"github.com/zajicekn/Pokedex-CLI/pokestructs"
)

var globalConfig pokestructs.Config

type cliCommand struct {
	name        string
	description string
	callback    func(area_name string) error
}

func commandExit(area_name string) error {
	if area_name == "" || area_name != "" {
		fmt.Println("Bye now!")
		os.Exit(0)
		return nil
	}
	return nil
}

func commandMap(area_name string) error {
	if area_name == "" {
		pokeAPI.Map(&globalConfig)
		return nil
	}
	return errors.New("map shouldn't have an argument")
}

func commandMapb(area_name string) error {
	if area_name == "" {
		pokeAPI.Mapb(&globalConfig)
		return nil
	}
	return errors.New("mapb shouldn't have an argument")
}

func commandExplore(area_name string) error {
	if area_name != "" {
		pokeAPI.Explore(&globalConfig, area_name)
		return nil
	}
	return errors.New("explore should have an argument")
}

func commandCatch(area_name string) error {
	if area_name != "" {
		pokeAPI.Catch(&globalConfig, area_name)
		return nil
	}
	return errors.New("catch should have an argument")
}

func commandInspect(area_name string) error {
	if area_name != "" {
		pokeAPI.Inspect(area_name)
		return nil
	}
	return errors.New("inspect should have an argument")
}

func commandMoves(area_name string) error {
	if area_name != "" {
		pokeAPI.Moves(area_name)
		return nil
	}
	return errors.New("moves should have an argument")
}

func commandPokedex(area_name string) error {
	if area_name == "" {
		pokeAPI.Pokedex()
		return nil
	}
	return errors.New("pokedex shouldn't have an argument")
}

func commandHelp(area_name string) error {
	if area_name == "" {
		color.Set(color.FgHiYellow)
		fmt.Println("\nWelcome to the Pokedex help menu!")
		fmt.Print("Commands are listed below\n\n")
		color.Unset()

		display := map[string]cliCommand{
			"help": {
				name:        color.HiYellowString("help"),
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"exit": {
				name:        color.HiYellowString("exit"),
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"map": {
				name:        color.HiYellowString("map"),
				description: "Displays the location names of 20 different areas in the Pokemon world",
				callback:    commandMap,
			},
			"mapb": {
				name:        color.HiYellowString("mapb"),
				description: "Displays the previous names of the 20 different areas in the Pokemon world",
				callback:    commandMapb,
			},
			"explore": {
				name:        color.HiYellowString("explore"),
				description: "Explores an area and presents the Pokemon in that area",
				callback:    commandExplore,
			},
			"catch": {
				name:        color.HiYellowString("catch"),
				description: "Catches a pokemon",
				callback:    commandCatch,
			},
			"inspect": {
				name:        color.HiYellowString("inspect"),
				description: "Inspects a pokemon and its stats",
				callback:    commandInspect,
			},
			"moves": {
				name:        color.HiYellowString("moves"),
				description: "Displays the moves for a certain pokemon",
				callback:    commandMoves,
			},
			"pokedex": {
				name:        color.HiYellowString("pokedex"),
				description: "Displays the names of pokemon in your pokedex",
				callback:    commandPokedex,
			},
		}
		for _, v := range display {
			fmt.Printf("%s: %s\n", v.name, v.description)
		}
		fmt.Print("\n")
		return nil
	}
	return errors.New("help shouldn't have an argument")
}

func main() {
	pokecache.InitGlobalCache()
	pokedex.InitGlobalDex()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nWelcome to the Pokedex-CLI! (type 'help' for more information)\n\n")

	for {
		color.Set(color.FgCyan)
		fmt.Print("Pokedex > ")
		color.Unset()

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		// Trim whitespace from input
		input = strings.TrimSpace(input)

		// Split the input from command to argument
		splitArgs := strings.Split(input, " ")

		switch {
		case splitArgs[0] == "help":
			commandHelp("")
		case splitArgs[0] == "exit":
			commandExit("")
		case splitArgs[0] == "map":
			commandMap("")
		case splitArgs[0] == "mapb":
			commandMapb("")
		case splitArgs[0] == "explore":
			if len(splitArgs) < 2 {
				log.Print("Need to have an argument for 'explore'")
			}
			commandExplore(splitArgs[1])
		case splitArgs[0] == "catch":
			if len(splitArgs) < 2 {
				log.Print("Need to have an argument for 'catch'")
			}
			commandCatch(strings.ToLower(splitArgs[1]))
		case splitArgs[0] == "moves":
			if len(splitArgs) < 2 {
				log.Print("Need to have an argument for 'moves'")
			}
			commandMoves(strings.ToLower(splitArgs[1]))
		case splitArgs[0] == "inspect":
			if len(splitArgs) < 2 {
				log.Print("Need to have an argument for 'inspect'")
			}
			commandInspect(strings.ToLower(splitArgs[1]))
		case splitArgs[0] == "pokedex":
			commandPokedex("")
		default:
			fmt.Printf("No such command '%s'. Please try again\n", input)
		}
	}
}
