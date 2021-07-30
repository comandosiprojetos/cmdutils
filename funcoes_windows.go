package cmdutils

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Recebe uma string com o comando a ser executado e retorna o stdout do comando e um tipo error
// Utilizado apenas no windows
func ExecutaComandoTerminalWindows(command string) (string, error) {
	cmd := exec.Command("cmd", "/C", command)

	stdOut, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("[ERRO] %s", err.Error())
	}

	return string(stdOut), nil
}

// Recebe uma string com o nome de um executável e retorna uma o string com o
// path absoluto de um executável no windows e um tipo error
// comando semelhante ao wish do linux
// o nome do executável informado precisa estar no path global do windows
func RetornaPathCompletoExecutavelWindows(nomeExe string) (string, error) {
	comando := fmt.Sprintf("%s in (%s) do @echo.   %s", "for %i", nomeExe, "%~$PATH:i")

	cmd := exec.Command("cmd", "/C", comando)

	stdOut, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	pathRetornado := strings.TrimSpace(strings.ReplaceAll(string(stdOut), `\`, `\\`))

	return pathRetornado, nil
}

// Retorna uma string com a versão do windows lendo a informação do registro do windows
// e retorna um tipo error
// exemplo retorno: Windows 10 Pro
func RetornaVersaoSistemaOperacional() (string, error) {
	registroCurrentVersion, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("Não foi possível abrir o registro para retornar a versão do windows. Erro: %s.",
			err.Error())
	}
	defer registroCurrentVersion.Close()

	atributoProductName, _, err := registroCurrentVersion.GetStringValue("ProductName")
	if err != nil {
		return "", fmt.Errorf("Não foi possível obter o atributo 'ProductName' do registro do windows. Erro: %s.",
			err.Error())
	}

	if atributoProductName == "" {
		return "", fmt.Errorf("O atributo 'ProductName' do registro do windows não foi localizado.")
	}

	return atributoProductName, nil
}

// Retorna uma string com o modelo do processador lendo a informação do registro do windows
// e retorna um tipo error
// exemplo retorno: Intel(R) Core(TM) i7-4770 CPU @ 3.40GHz
func RetornaModeloProcessadorWindows() (string, error) {
	registroCentralProcessor, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("Não foi possível abrir o registro para retornar o modelo do processador no windows. Erro: %s.",
			err.Error())
	}
	defer registroCentralProcessor.Close()

	atributoProcessorNameString, _, err := registroCentralProcessor.GetStringValue("ProcessorNameString")
	if err != nil {
		return "", fmt.Errorf("Não foi possível obter o atributo 'ProcessorNameString' do registro do windows. Erro: %s.",
			err.Error())
	}

	if atributoProcessorNameString == "" {
		return "", fmt.Errorf("O atributo 'ProcessorNameString' do registro do windows não foi localizado.")
	}

	return atributoProcessorNameString, nil
}
