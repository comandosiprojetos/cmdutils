package cmdutils

import (
	"fmt"
	"strings"
	"bytes"
)

// Retorna uma string com a versão da distro linux e um tipo error
func RetornaVersaoSistemaOperacional() (string, error) {
	comando := `awk -F= '$1=="PRETTY_NAME" { print $2 ;}' /etc/os-release`

	retornoComando, errComando := cmdutils.ExecutaComandoTerminalLinux(comando)
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
func RetornaModeloProcessadorLinux() (string, error) {
	comando := `cat /proc/cpuinfo | grep 'model name' | uniq`

	retornoComando, errComando := cmdutils.ExecutaComandoTerminalLinux(comando)
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

// Recebe uma string com o comando a ser executado e retorna o stdout do comando e um tipo error
func ExecutaComandoTerminalLinux(tipoComando string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", "-c", tipoComando)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("[ERRO] %s", err.Error())
	}

	return stdout.String(), nil
}

// Retorna uma string informando o chassis sob o qual a distro linux está sendo executada 
// e também retorna um tipo error
func ChassisMaquinaLinux() (string, error) {
	cmd := "hostnamectl | grep Chassis:"
	retornoComando, errCmd := cmdutils.ExecutaComandoTerminalLinux(cmd)
	if errCmd != nil {
		return "", errCmd
	}

	retornoComNomeChassiRemovido := strings.ReplaceAll(retornoComando, "Chassis:", "")
	retornoComEspacosRemovido := strings.TrimSpace(retornoComNomeChassiRemovido)

	return retornoComEspacosRemovido, nil
}

// Retorna uma string informando a versão do kernel da distro linux instalada
// e também retorna um tipo error
func VersaoKernelMaquinaLinux() (string, error) {
	cmd := "hostnamectl | grep Kernel:"
	retornoComando, errCmd := cmdutils.ExecutaComandoTerminalLinux(cmd)
	if errCmd != nil {
		return "", errCmd
	}

	retornoComNomeKernelRemovido := strings.ReplaceAll(retornoComando, "Kernel:", "")
	retornoComEspacosRemovido := strings.TrimSpace(retornoComNomeKernelRemovido)

	return retornoComEspacosRemovido, nil
}
