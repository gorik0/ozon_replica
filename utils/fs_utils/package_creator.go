package fs_utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func CreatePkgs(root string, pkgs string) error {
	for _, pkg := range strings.Split(pkgs, ",") {

		err := os.Mkdir(filepath.Join(root, pkg), 0755)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				println("let it be ... ")

			} else {
				return err

			}
		}

	}
	return nil
}
