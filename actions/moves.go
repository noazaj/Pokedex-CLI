package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/zajicekn/Pokedex-CLI/pokedex"
)

func GetMoves(name string) ([]string, error) {
	dex := pokedex.GetGlobalDex()

	_, ok := dex.Data[name]
	if !ok {
		return nil, errors.New("you have not caught that pokemon")
	}

	url := fmt.Sprintf("https://cs361-get-json-key-micro-0d1965251b55.herokuapp.com/pokemon/%s", name)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var moves []string
	err = json.Unmarshal(respData, &moves)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return moves, nil
}
