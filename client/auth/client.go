package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"ozon_replic/internal/pkg/auth/delivery/grpc/gen"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, err := grpc.NewClient("localhost:46783", opts...)
	if err != nil {
		panic("Failed to create client" + err.Error())
	}

	cli := gen.NewAuthClient(client)

	whatToCheck := os.Getenv("WHAT_TO_CHECK")

	switch whatToCheck {
	case "in":
		{
			trySignIn(err, cli)

		}
	case "up":
		{
			trySignUp(err, cli)

		}
	case "check":
		{
			tryCheckAuth(err, cli)

		}
	default:
		println("There is no anuy choise to chek!!!!")

	}

}

func tryCheckAuth(err error, cli gen.AuthClient) {
	chekReq := gen.CheckAuthRequst{
		ID: "b49d0106-0937-427f-bf57-e44591453f35",
	}
	auth, err := cli.CheckAuth(context.Background(), &chekReq)
	if err != nil {
		panic("Failed to check auth" + err.Error())

	}
	log.Printf("Checking auth is success, %v", auth)
}

func trySignUp(err error, cli gen.AuthClient) {
	siUp := gen.SignUpRequest{
		Login:    "gorikos",
		Password: "121341341234124",
		Phone:    "+14155552671",
	}
	up, err := cli.SignUp(context.Background(), &siUp)
	if err != nil {
		println(up)
		panic(err)
	}
	println("DONE UP!!! .... ", up)
}

func trySignIn(err error, cli gen.AuthClient) {
	siIn := gen.SignInRequest{
		Login:    "gorikos",
		Password: "121341341234124",
	}
	in, err := cli.SignIn(context.Background(), &siIn)
	if err != nil {
		panic(err)
	}
	log.Printf("DONE IN!!! .... %v", in)

}
