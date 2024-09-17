package main

import (
	"os"
	"ozon_replic/utils/fs_utils"
	"path/filepath"
)

func main() {

	pwd, _ := os.Getwd()
	err := fs_utils.CreatePkgs(filepath.Join(pwd, "internal/pkg/profile"), "delivery,mocks,repo,usecase")
	if err != nil {
		panic(err)
	}

}
