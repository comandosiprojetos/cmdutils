package cmdutils

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

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
func RetornaModeloProcessador() (string, error) {
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
