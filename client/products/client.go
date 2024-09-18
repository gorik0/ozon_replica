package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	gen2 "ozon_replic/internal/pkg/products/delivery/grpc/gen"
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

	cli := gen2.NewProductsClient(client)
	_ = cli

	switch whatToCheck {
	case "get_one":
		{
			tryGetProduct(err, cli)

		}
	case "get_many":
		{
			tryGetProducts(err, cli)

		}
	case "get_cat":
		{
			tryGetCategory(err, cli)

		}
	default:
		println("There is no anuy choise to chek!!!!")

	}

}

func tryGetCategory(err error, cli gen2.ProductsClient) {
	in := gen2.CategoryRequest{
		Id:     21,
		Paging: 1,
		Count:  30,
	}
	category, err := cli.GetCategory(context.Background(), &in)
	if err != nil {
		panic("Failed to get category" + err.Error())
	}
	log.Printf("Category: %v", category)
	for _, product := range category.Products {
		println(":::::::::")
		log.Printf("Product: %v", product.Name)

	}
}

func tryGetProducts(err error, cli gen2.ProductsClient) {
	in := gen2.ProductsRequest{
		Paging: 1,
		Count:  20,
	}
	products, err := cli.GetProducts(context.Background(), &in)
	if err != nil {
		panic("Failed to get products: " + err.Error())
	}
	log.Printf("Got %d products", len(products.Products))
	for _, product := range products.Products {
		log.Printf("Product: %s", product.Name)

	}
}

func tryGetProduct(err error, cli gen2.ProductsClient) {
	in := gen2.ProductRequest{
		Id: "8472b67f-2d05-449f-8138-2b87f87eea24",
	}
	product, err := cli.GetProduct(context.Background(), &in)
	if err != nil {
		panic("Failed to get product" + err.Error())
	}
	log.Printf("Success on getting product :: %v", product)
}
