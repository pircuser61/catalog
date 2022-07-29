package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
)

func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewCatalogClient(conn)
	ctx := context.Background()

	var line string
	in := bufio.NewScanner(os.Stdin)

	defer func() { fmt.Println("...press enter"); in.Scan() }()
	for {
		fmt.Print("\n>")

		if !in.Scan() {
			fmt.Println("Scan error")
			continue
		}
		line = in.Text()
		cmd := strings.Split(line, " ")[0]
		switch cmd {
		case "quit":
			fallthrough
		case "q":
			return
		case "list":
			response, err := client.GoodList(ctx, &emptypb.Empty{})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "add":
			params := strings.Split(line, " ")
			if len(params) != 4 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			request := pb.GoodCreateRequest{Name: params[1], UnitOfMeasure: params[2], Country: params[3]}
			response, err := client.GoodCreate(ctx, &request)
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "update":
			params := strings.Split(line, " ")
			if len(params) != 5 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			code, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fmt.Println("<code> must be a number")
				continue
			}
			request := pb.GoodUpdateRequest{Code: code, Name: params[2], UnitOfMeasure: params[3], Country: params[4]}
			response, err := client.GoodUpdate(ctx, &request)
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "get":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			code, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fmt.Println("<code> must be a number")
				continue
			}
			response, err := client.GoodGet(ctx, &pb.GoodGetRequest{Code: code})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "delete":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			code, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fmt.Println("<code> must be a number")
				continue
			}
			response, err := client.GoodDelete(ctx, &pb.GoodDeleteRequest{Code: code})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}

		default:
			fmt.Printf("Unknown command <%s>\n", line)
		}
	}
}
