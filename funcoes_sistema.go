package cmdutils

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
)

// Retorna uma string com o tipo de barras utilizado no linux ou no windows
func RetornaBarrasOs() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}

	return "/"
}

// Retorna um valor boelando informando se existe ou não conexão com a internet
func ConexaoWebAtiva() bool {
	_, err := http.Get("https://www.google.com.br/")
	if err != nil {
		return false
	}

	return true
}

// Valida a utilização de uma única instância abrindo um porta tcp para teste e retorna um boleano
func ValidarInstanciaUnica(portaTcp string) bool {
	if _, err := net.Listen("tcp", fmt.Sprintf(":%s", portaTcp)); err != nil {
		return true
	}

	return false
}
