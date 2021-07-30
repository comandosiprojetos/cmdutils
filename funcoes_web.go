package cmdutils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
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

type IP struct {
	Query string
}

// Retorna uma string com o ip público da rede e um tipo error
// Função utiliza a api http://ip-api.com/json/
func RetornaIpPublico() (string, error) {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query, nil
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

// Recebe uma string com uma url e uma string com o nome do parâmetro a ser extraído
// e retorna uma string com esse parâmetro extraído da url e um tipo error
func RetornaParametroPresenteEmUmaUrlInformada(urlInformada, nomeParametro string) (error, string) {
	u, errParse := url.Parse(urlInformada)
	if errParse != nil {
		return errParse, ""
	}

	urlRaw, errParseQuery := url.ParseQuery(u.RawQuery)
	if errParseQuery != nil {
		return errParseQuery, ""
	}

	if len(urlRaw[nomeParametro]) == 0 {
		return fmt.Errorf("O parâmetro informado '%s' não pode ser localizado na url '%s' informada. ", nomeParametro, urlInformada), ""
	}

	return nil, urlRaw[nomeParametro][0]
}

// Recebe um inteiro com a porta tcp a ser verificada e checa se a mesma está ou não aberta
// retorna um tipo booleano e um tipo error
func PortaTcpEstaAberta(host string, portaTCP int) (bool, error) {
	timeOutConexao := time.Second
	conexao, errConexao := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprint(portaTCP)), timeOutConexao)
	if errConexao != nil {
		return false, nil
	}

	if conexao != nil {
		defer conexao.Close()
		return true, nil
	}

	return false, nil
}

// Retorna uma instância do tipo net.Listener e um tipo error
func RetornaInstanciaNetConn(host string, portaTCP int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", host, portaTCP))
}
