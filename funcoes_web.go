package cmdutils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

// Retorna um valor booleano informando se existe ou não conexão com a internet
func ConexaoWebAtiva() bool {
	_, err := http.Get("https://www.google.com.br/")
	if err != nil {
		return false
	}

	return true
}

// Valida a utilização de uma única instância abrindo um porta tcp para teste e retorna um booleano
func ValidarInstanciaUnica(portaTcp string) bool {
	if _, err := net.Listen("tcp", fmt.Sprintf(":%s", portaTcp)); err != nil {
		return true
	}

	return false
}

// Recebe um string com o caminho do arquivo e outra string com a url de download
// retorna um tipo error
func DownloadArquivo(localArquivo string, urlDownload string) error {
	resp, err := http.Get(urlDownload)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(localArquivo)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

// Retorna uma string com o ip local da máquina e um tipo error
func RetornaIpLocalMaquina() (string, error) {
	con, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer con.Close()

	localAddr := con.LocalAddr().(*net.UDPAddr)

	return fmt.Sprint(localAddr.IP), nil
}

// Recebe uma string com uma url e retorna um valor booleano para informar se a url é válida ou não
func PingUrl(url string) bool {
	resp, netErrors := http.Get(url)
	if netErrors != nil {
		return false
	}

	defer resp.Body.Close()

	return true
}
