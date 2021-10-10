package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const defaultPort = "4041"

func main() {

	fmt.Println("G Client")

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

	c := grpc.NewGServiceClient(cc)

	// create G
	fmt.Println("Creating the g")
	g := &grpc.G{
		Pid:         "G01",
		Name:        "G",
		Power:       "Fire",
		Description: "Fluffy",
	}
	createGRes, err := c.CreateG(context.Background(), &grpc.CreateGRequest{G: g})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("G has been created: %v", createGRes)
	gID := createGRes.GetG().GetId()

	// read G
	fmt.Println("Reading the g")
	readGReq := &grpc.ReadGRequest{Pid: gID}
	readGRes, readGErr := c.ReadG(context.Background(), readGReq)
	if readGErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readGErr)
	}

	fmt.Printf("G was read: %v \n", readGRes)

	// update G
	newG := &grpc.G{
		Id:          gID,
		Pid:         "G01",
		Name:        "G",
		Power:       "Fire Fire Fire",
		Description: "Fluffy",
	}
	updateRes, updateErr := c.UpdateG(context.Background(), &grpc.UpdateGRequest{Pokemon: newG})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("G was updated: %v\n", updateRes)

	// delete G
	deleteRes, deleteErr := c.DeleteG(context.Background(), &grpc.DeleteGRequest{Pid: gID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("G was deleted: %v \n", deleteRes)

	// list G

	stream, err := c.ListG(context.Background(), &grpc.ListGRequest{})
	if err != nil {
		log.Fatalf("error while calling ListG RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetG())
	}
}
