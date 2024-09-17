package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"ozon_replic/internal/pkg/auth/delivery/grpc/gen"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, err := grpc.NewClient("localhost:35591", opts...)
	if err != nil {
		panic("Failed to create client" + err.Error())
	}

	cli := gen.NewAuthClient(client)

	siIn := gen.SignInRequest{
		Login:    "gorikos",
		Password: "121341341234124",
	}
	up, err := cli.SignIn(context.Background(), &siIn)
	if err != nil {
		panic(err)
	}
	//siIn := gen.SignUpRequest{
	//	Login:    "gorikos",
	//	Password: "121341341234124",
	//	Phone:    "+14155552671",
	//}
	//up, err := cli.SignUp(context.Background(), &siIn)
	//if err != nil {
	//	panic(err)
	//}
	println("DONE!!! .... ", up)

}
