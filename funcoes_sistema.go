package cmdutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/denisbrodbeck/machineid"
)

// Retorna uma string com o tipo de barras utilizado no linux ou no windows
func RetornaBarrasOs() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}

	return "/"
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

// Retorna uma string com o diretorio do próprio executável em execução e um tipo error
func RetornaDiretorioAplicacao() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	exPath := fmt.Sprint(filepath.Dir(ex), RetornaBarrasOs())

	return exPath, nil
}

// Retorna uma string com o caminho exato do executável em execução e um tipo error
func RetornaCaminhoAbsolutoAplicacao() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return ex, nil
}

// Retorna uma string com o nome do executável em execução e um tipo error
func RetornaNomeDoExeDaAplicacao() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Base(ex), nil
}

// Recebe um []byte com uma chave de criptografia e uma string com a mensagem a ser criptografada
// Retorna uma string com a mensagem criptografada e um tipo error
func CriptografarTextoUtilizandoAES(chave []byte, mensagem string) (mensagemCriptografada string, err error) {
	mensagemOriginal := []byte(mensagem)

	block, err := aes.NewCipher(chave)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(mensagemOriginal))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], mensagemOriginal)

	mensagemCriptografada = base64.URLEncoding.EncodeToString(cipherText)

	return
}

// Recebe um []byte com uma chave de criptografia e uma string com a mensagem criptografada
// Retorna uma string com a mensagem descriptografada e um tipo error
func DescriptografarTextoUtilizandoAES(chave []byte, mensagemCriptografada string) (mensagemDecodificada string, err error) {
	textoCriptografado, err := base64.URLEncoding.DecodeString(mensagemCriptografada)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(chave)
	if err != nil {
		return
	}

	if len(textoCriptografado) < aes.BlockSize {
		err = errors.New("O tamanho do bloco de texto cifrado é muito curto!")
		return
	}

	iv := textoCriptografado[:aes.BlockSize]
	textoCriptografado = textoCriptografado[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(textoCriptografado, textoCriptografado)

	mensagemDecodificada = string(textoCriptografado)

	return
}

// Retorna uma string com um código uuid aleatório
func RetornarUUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

// Recebe uma string com o local do arquivo e retorna uma string com o md5 do arquivo calculado
// e um tipo error
func CalcularMD5Arquivo(localArquivo string) (string, error) {
	f, err := os.Open(localArquivo)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Recebe uma string com o caminho de um diretório e apaga recursivamente o diretório
// e retorna um tipo error
func RemoverPasta(caminhoDiretorio string) error {
	listaArquivos, errGlob := filepath.Glob(filepath.Join(caminhoDiretorio, "*"))
	if errGlob != nil {
		return errGlob
	}

	for _, arquivo := range listaArquivos {
		errRemoveAll := os.RemoveAll(arquivo)
		if errRemoveAll != nil {
			return errRemoveAll
		}
	}

	errRemove := os.Remove(caminhoDiretorio)
	if errRemove != nil {
		return errRemove
	}

	return nil
}

// Recebe uma string com o local e nome da pasta a ser criada e ignora a criação caso a pasta já exista
// e retorna um error
func CriarPastaIgnorandoCasoJaExista(nomePasta string) error {
	if _, errStat := os.Stat(nomePasta); !os.IsNotExist(errStat) {
		return nil
	} else {
		os.MkdirAll(nomePasta, os.ModePerm)

		if _, errStat := os.Stat(nomePasta); os.IsNotExist(errStat) {
			return fmt.Errorf("Não foi possível criar a pasta '%s'. Erro original: %s", nomePasta, errStat.Error())
		}
	}

	return nil
}

// Recebe uma string com o local e nome da pasta a ser criada e retorna um erro caso a pasta já exista
// retorna um error
func CriarPastaValidandoSeAPastaExiste(nomePasta string) error {
	if _, errStat := os.Stat(nomePasta); !os.IsNotExist(errStat) {
		return fmt.Errorf("Não foi possível criar a pasta '%s' Pasta já existe. Erro original: %s", nomePasta, errStat.Error())
	} else {
		os.MkdirAll(nomePasta, os.ModePerm)

		if _, errStat := os.Stat(nomePasta); os.IsNotExist(errStat) {
			return fmt.Errorf("Não foi possível criar a pasta '%s'. Erro original: %s", nomePasta, errStat.Error())
		}
	}

	return nil
}

// Realiza a compressão de um arquivo utilizando o 7Zip
// No windows irá utilizar as dlls presentes na pasta da aplicação
// equanto que no linux será necessário ter um pacote do 7Zip instalado no sistema operacional
// Recebe uma string com o nome do arquivo já com a extenção .7z e outra string com o nome original do arquivo
func ComprimirArquivoCom7Zip(nomeArquivo7zip string, nomeArquivo string) error {
	var local7Zip string

	local7Zip = "7z"

	if runtime.GOOS == "windows" {
		pathAbsolutoAplicacao, errRetornaDiretorio := RetornaDiretorioAplicacao()
		if errRetornaDiretorio != nil {
			return errRetornaDiretorio
		}

		local7Zip = fmt.Sprint(pathAbsolutoAplicacao, "7z")
	}

	cmd := exec.Command(local7Zip, "a", nomeArquivo7zip, nomeArquivo)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Realiza a descompressão de um arquivo utilizando o 7Zip
// No windows irá utilizar as dlls presentes na pasta da aplicação
// equanto que no linux será necessário ter um pacote do 7Zip instalado no sistema operacional
// Recebe uma string com o local exato do arquivo, uma string com o diretório para extração
// e uma string com o local do executável
func DescomprimirArquivoCom7Zip(localArquivo string, pastaDestino string) error {
	local7Zip := ""

	if runtime.GOOS == "windows" {
		pathAbsolutoAplicacao, errRetornaDiretorio := RetornaDiretorioAplicacao()
		if errRetornaDiretorio != nil {
			return errRetornaDiretorio
		}

		local7Zip = fmt.Sprint(pathAbsolutoAplicacao, RetornaBarrasOs(), "7z")
	}

	cmd := exec.Command(fmt.Sprintf("%s%s", local7Zip, "7z"), "x", localArquivo, fmt.Sprint("-o", pastaDestino))

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Recebe uma string com o arquivo compactado com 7Zip e retorna um booleano
// informando se o arquivo está ou não integro após sua compressão
func TesteCompressaoRetornouOK(nomeArquivo7zip string) bool {
	local7Zip := "7z"

	if runtime.GOOS == "windows" {
		pathAbsolutoAplicacao, errRetornaDiretorio := RetornaDiretorioAplicacao()
		if errRetornaDiretorio != nil {
			return false
		}

		local7Zip = fmt.Sprint(pathAbsolutoAplicacao, RetornaBarrasOs(), "7z")
	}

	cmd := exec.Command(local7Zip, "t", nomeArquivo7zip)

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

// Retorna se um arquivo existe ou não.
// Recebe uma string com o caminho do arquivo e retorna um valor boleando
func ArquivoExiste(caminhoArquivo string) bool {
	if _, err := os.Stat(caminhoArquivo); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// Recebe duas strings a primeira com o nome antigo do arquivo
// e a segunda com novo nome do arquivo
func RenomearArquivo(nomeAntigo string, novoNome string) error {
	return os.Rename(nomeAntigo, novoNome)
}

// Recebe uma string com um valor e retorna uma string apenas com números
// e letras e um valor do tipo error
func RemoverCaracteresEspeciais(valor string) (string, error) {
	regex, errRegex := regexp.Compile("[^a-zA-Z0-9 ]+")
	if errRegex != nil {
		return "", errRegex
	}

	return regex.ReplaceAllString(valor, ""), nil
}

// Recebe duas datas do tipo time e retorna a diferença entre elas em uma string
func ImprimirDiferencaEntreDuasDatas(dataHora1 time.Time, dataHora2 time.Time) string {
	hs := dataHora1.Sub(dataHora2).Hours()

	hs, mf := math.Modf(hs)
	ms := mf * 60

	ms, sf := math.Modf(ms)
	ss := sf * 60

	return fmt.Sprint(math.Abs(hs), " horas, ", math.Abs(ms), " minutos, ", math.Abs(ss), " segundos ")
}

// Recebe uma string com o caminho do arquivo e o remove caso existe e retorna um tipo errror
func RemoverArquivo(localArquivo string) error {
	if _, err := os.Stat(localArquivo); !os.IsNotExist(err) {
		err := os.Remove(localArquivo)
		if err != nil {
			return fmt.Errorf("Ocorreu um erro ao apagar o arquivo '%s'. Erro: %s", localArquivo, err.Error())
		}
	}

	return nil
}

// Recebe uma string com o caminho do arquivo retorna um tipo int64 com o tamanho do arquivo
// e um tipo errror
func RetornaTamanhoArquivo(localArquivo string) (int64, error) {
	fi, err := os.Stat(localArquivo)
	if err != nil {
		return -1, err
	}

	size := fi.Size()

	return size, nil
}

// Recebe um float64 com o valor em bytes e formata de
// acordo com o tamanho retornando uma string com o valor em B, KB, MB, GB ou TB
func FormatarValorEmBytes(tamanhoEmBytes float64) string {
	var sufixo [5]string

	sufixo[0] = "B"
	sufixo[1] = "KB"
	sufixo[2] = "MB"
	sufixo[3] = "GB"
	sufixo[4] = "TB"

	base := math.Log(tamanhoEmBytes) / math.Log(1024)
	tamanhoCalculado := ArredondarValor(math.Pow(1024, base-math.Floor(base)), .5, 2)
	sufixoRetornado := sufixo[int(math.Floor(base))]

	return fmt.Sprint(strconv.FormatFloat(tamanhoCalculado, 'f', -1, 64) + " " + string(sufixoRetornado))
}

// Recebe um tipo float64 com o valor a ser arredondado, um valor do tipo float64
// com o digito que deve ser considerado no arredondamento e um valor do tipo int
// indicando em quantas casas decimais o valor será arredondado
func ArredondarValor(valorOriginal float64, arredondarEm float64, casasDecimais int) float64 {
	var round float64

	pow := math.Pow(10, float64(casasDecimais))
	digito := pow * valorOriginal
	_, divisao := math.Modf(digito)

	if divisao >= arredondarEm {
		round = math.Ceil(digito)
	} else {
		round = math.Floor(digito)
	}

	return round / pow
}

// Recebe o caminho de um diretório e retorna uma lista do tipo string
// com os arquivos presentes nesse diretório
func RetornaListaDeArquivosDeUmDiretorio(diretorio string) ([]string, error) {
	var files []string

	err := filepath.Walk(diretorio, func(path string, info os.FileInfo, err error) error {
		if path != diretorio {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Ocorreu um erro ao listar os arquivos do diretório. Erro: %s", err.Error())
	}

	return files, nil
}

// Recebe uma string com data e hora no padrão 2021-11-11 23:59:59
// e retorna uma string no padrão 2021-11-11T23:59:59.000-0200
func FormataDataPadraoISO8601(DataHora string) string {
	return fmt.Sprintf("%s%s", strings.Replace(DataHora, " ", "T", 1), ".000-0200")
}

// Recebe um slice em interface e um indice do tipo int e remove
// o item desse slice de acordo com o indice informado
func RemoverItemSlice(Slice interface{}, Indice int) {
	vField := reflect.ValueOf(Slice)

	value := vField.Elem()

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		result := reflect.AppendSlice(value.Slice(0, Indice), value.Slice(Indice+1, value.Len()))

		value.Set(result)
	}
}

// Retorna uma string com a data e hora no padrão dd-mm-yyyy hh:mm:ss
func RetornaDataHora() string {
	now := time.Now()

	return now.Format("02-01-2006 15:04:05")
}

// Retorna uam srinf com o md5 do próprio exeutável da aplicação em execução
func RetornarMD5Aplicacao() (string, error) {
	caminhoAplicacao, errRetornaCaminho := RetornaCaminhoAbsolutoAplicacao()
	if errRetornaCaminho != nil {
		return "", errRetornaCaminho
	}

	if _, err := os.Stat(caminhoAplicacao); os.IsNotExist(err) {
		return "", nil
	}

	return CalcularMD5Arquivo(caminhoAplicacao)
}

// Recebe uma string com um valor textual e um int com a quantidade
// de digitos numéricos que esse valor textual deve conter e retorna um valor booleano
// para confirmar ou negar se a string passada possui a quantidade de digitos esperada ou não
func StringPossuiQuantidadeDeDigitosNumericosCorreta(valor string, quantidadeDigitos int) bool {
	var valorCalculado string

	validarDigitosNumericos := regexp.MustCompile("[0-9]+")
	valorExtraido := validarDigitosNumericos.FindAllString(valor, -1)

	for indice := range valorExtraido {
		valorCalculado = valorCalculado + valorExtraido[indice]
	}

	if len(valorCalculado) < quantidadeDigitos {
		return false
	}

	return true
}

// Recebe um valor do tipo []byte e retorna uma string com o valor convertido
func BytesToString(valorEmBytes []byte) string {
	return string(valorEmBytes[:])
}

// Recebe o nome de uma variável de ambiente reseta a mesma e retorna um tipo error
func ResetarVariaveisAmbiente(nomeVariavel string) error {
	err := os.Unsetenv(nomeVariavel)
	if err != nil {
		return err
	}

	return nil
}

// Recebe uma string com uma chave criptografica e retorna uma string com o id 
// único gerado e um tipo error
func RetornaIdUnicoMaquina(chaveCriptografica string) (string, error) {
	id, err := machineid.ProtectedID(chaveCriptografica)
	if err != nil {
		return "", err
	}

	return id, nil
}
