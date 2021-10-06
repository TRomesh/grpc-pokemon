package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pokemonpc "github.com/TRomesh/grpc-pokemon/pokemon"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const defaultPort = "4041"

func main() {

	fmt.Println("Pokemon Client")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:4041", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close() // Maybe this should be in a separate function and the error handled?

	c := pokemonpc.NewPokemonServiceClient(cc)

	// create Pokemon
	fmt.Println("Creating the pokemon")
	pokemon := &pokemonpc.Pokemon{
		Pid:         "Poke01",
		Name:        "Pikachu",
		Power:       "Fire",
		Description: "Fluffy",
	}
	createPokemonRes, err := c.CreatePokemon(context.Background(), &pokemonpc.CreatePokemonRequest{Pokemon: pokemon})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Pokemon has been created: %v", createPokemonRes)
	pokemonID := createPokemonRes.GetPokemon().GetId()

	// read Pokemon
	fmt.Println("Reading the pokemon")
	readPokemonReq := &pokemonpc.ReadPokemonRequest{Pid: pokemonID}
	readPokemonRes, readPokemonErr := c.ReadPokemon(context.Background(), readPokemonReq)
	if readPokemonErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readPokemonErr)
	}

	fmt.Printf("Pokemon was read: %v \n", readPokemonRes)

	// update Pokemon
	newPokemon := &pokemonpc.Pokemon{
		Id:          pokemonID,
		Pid:         "Poke01",
		Name:        "Pikachu",
		Power:       "Fire Fire Fire",
		Description: "Fluffy",
	}
	updateRes, updateErr := c.UpdatePokemon(context.Background(), &pokemonpc.UpdatePokemonRequest{Pokemon: newPokemon})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Pokemon was updated: %v\n", updateRes)

	// delete Pokemon
	deleteRes, deleteErr := c.DeletePokemon(context.Background(), &pokemonpc.DeletePokemonRequest{Pid: pokemonID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Pokemon was deleted: %v \n", deleteRes)

	// list Pokemons

	stream, err := c.ListPokemon(context.Background(), &pokemonpc.ListPokemonRequest{})
	if err != nil {
		log.Fatalf("error while calling ListPokemon RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetPokemon())
	}
}
