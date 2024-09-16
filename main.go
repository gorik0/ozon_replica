package main

import (
	"os"
	"ozon_replic/utils/fs_utils"
	"path/filepath"
)

func main() {

	pwd, _ := os.Getwd()
	err := fs_utils.CreatePkgs(filepath.Join(pwd), "metrics,postman,proto")
	if err != nil {
		panic(err)
	}

}
