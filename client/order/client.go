package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"ozon_replic/internal/pkg/order/delivery/grpc/gen"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	whatToCheck := os.Getenv("WHAT_TO_CHECK")
	port := os.Getenv("port")

	client, err := grpc.NewClient(":"+port, opts...)
	if err != nil {
		panic("Failed to create client" + err.Error())
	}

	cli := gen.NewOrderClient(client)
	_ = cli

	switch whatToCheck {
	case "create":
		{
			tryCreateOrder(err, cli)

		}
	//case "get_many":
	//	{
	//		tryGetProducts(err, cli)
	//
	//	}
	//case "get_cat":
	//	{
	//		tryGetCategory(err, cli)
	//
	//	}
	default:
		println("There is no anuy choise to chek!!!!")

	}

}

func tryCreateOrder(err error, cli gen.OrderClient) {
	in := gen.CreateOrderRequest{
		Id:            "b49d0106-0937-427f-bf57-e44591453f35",
		DeliveryDate:  "",
		DeliveryTime:  "",
		PromocodeName: "",
	}
	order, err := cli.CreateOrder(context.Background(), &in)
	if err != nil {
		panic("Failed to create order" + err.Error())
	}
	log.Printf("Successfully created order: %v", order)
}
