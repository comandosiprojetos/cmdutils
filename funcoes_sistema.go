package cmdutils

import "runtime"

func RetornaBarrasOs() string {
	var Barras string

	if runtime.GOOS == "windows" {
		Barras = "\\"
	} else {
		Barras = "/"
	}

	return Barras
}
