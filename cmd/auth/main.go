package main

import "os"

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {

	//	::: init CONFIG
	//	::: init  LOGGER
	//	::: inti  DB
	//	::: inti  REPO
	//	::: inti  USECASE
	//	::: inti  HANDLER
	//	::: init GRPC server

	//	::: go  metric server
	//	::: go  GRPC server
	//	::: gracefull STOP  GRPC

}
