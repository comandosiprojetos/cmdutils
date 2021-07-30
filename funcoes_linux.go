package cmdutils

import (
	"fmt"
	"strings"
)

// Retorna uma string com a versão da distro linux e um tipo error
func RetornaVersaoSistemaOperacional() (string, error) {
	comando := `awk -F= '$1=="PRETTY_NAME" { print $2 ;}' /etc/os-release`

	retornoComando, errComando := ExecutaComandoTerminalLinux(comando)
	if errComando != nil {
		return "", fmt.Errorf("Execução do comando '%s' falhou. Erro: %s.",
			comando, errComando.Error())
	}

	if retornoComando == "" {
		return "", fmt.Errorf("Não foi possível retornar a versão da distro linux.")
	}

	nomeSistema := strings.ReplaceAll(retornoComando, `"`, "")

	return nomeSistema, nil
}

// Retorna uma string com o modelo do processador da máquina e um tipo error
func RetornaModeloProcessador() (string, error) {
	comando := `cat /proc/cpuinfo | grep 'model name' | uniq`

	retornoComando, errComando := ExecutaComandoTerminalLinux(comando)
	if errComando != nil {
		return "", fmt.Errorf("Execução do comando '%s' falhou. Erro: %s.",
			comando, errComando.Error())
	}

	if retornoComando == "" {
		return "", fmt.Errorf("Não foi possível retornar o modelo do processador na distro linux.")
	}

	retornoSemModelName := strings.ReplaceAll(retornoComando, "model name", "")
	retornoSemDoisPontos := strings.ReplaceAll(retornoSemModelName, ":", "")
	retornoSemEspacos := strings.TrimSpace(retornoSemDoisPontos)

	return retornoSemEspacos, nil
}
