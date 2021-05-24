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
)

// Retorna uma string com o tipo de barras utilizado no linux ou no windows
func RetornaBarrasOs() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}

	return "/"
}

// Retorna uma o string com o path absoluto de um executável no windows comando semelhente ao wish do linux
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

// Retorna uma string com o path absoluto do executável compilado windows
func RetornaPathAbsolutoAplicacao() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	exPath := filepath.Dir(ex)

	return exPath, nil
}

func Encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)

	return
}

func Decrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)

	return
}

// Retorna uma string com um código uuid aleatório
func RetornarUUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

// Recebe uma string com o local do arquivo e retorna uma string com o md5 do arquivo calculado
func CalcularMD5Arquivo(localArquivo string) (string, error) {
	f, err := os.Open(localArquivo)
	if err != nil {
		return "", fmt.Errorf("Ocorreu um erro ao abrir o arquivo para calcular md5 do arquivo. Erro: %s", err.Error())
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("Ocorreu um erro ao calcular md5 do arquivo. Erro: %s", err.Error())
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Recebe uma string com o caminho de um diretório e apaga recursivamente o diretório
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

// Recebe uma string com o nome da pasta e a cria no local informado
func CriarPasta(nomePasta string, ignorarPastaCriada bool) error {
	if _, errStat := os.Stat(nomePasta); !os.IsNotExist(errStat) {
		if !ignorarPastaCriada {
			return fmt.Errorf("Não foi possível criar a pasta '%s' Pasta já existe. Erro original: %s", nomePasta, errStat.Error())
		}
	} else {
		os.MkdirAll(nomePasta, os.ModePerm)

		if _, errStat := os.Stat(nomePasta); os.IsNotExist(errStat) {
			return fmt.Errorf("Não foi possível criar a pasta '%s'. Erro original: %s", nomePasta, errStat.Error())
		}
	}

	return nil
}

func ComprimirArquivo(nomeArquivo7zip string, nomeArquivo string) error {
	var local7Zip string

	local7Zip = "7z"
	
	if runtime.GOOS == "windows" {
		pathAbsolutoAplicacao, _ := RetornaPathAbsolutoAplicacao()
		local7Zip = fmt.Sprint(pathAbsolutoAplicacao, RetornaBarrasOs(), "7z")
	}

	cmd := exec.Command(local7Zip, "a", nomeArquivo7zip, nomeArquivo)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func DescomprimirArquivo(localArquivo string, pastaDestino string, local7zip string) error {
	if runtime.GOOS == "linux" {
		local7zip = ""
	}

	cmd := exec.Command(fmt.Sprintf("%s%s", local7zip, "7z"), "x", localArquivo, fmt.Sprint("-o", pastaDestino))

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func TesteCompressaoRetornouOK(nomeArquivo7zip string) bool {
	var local7Zip string

	local7Zip = "7z"

	if runtime.GOOS == "windows" {
		pathAbsolutoAplicacao, _ := RetornaPathAbsolutoAplicacao()
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

func RenomearArquivo(nomeAntigo string, novoNome string) error {
	return os.Rename(nomeAntigo, novoNome)
}

func StripRegex(valor string) (string, error) {
	regex, errRegex := regexp.Compile("[^a-zA-Z0-9 ]+")
	if errRegex != nil {
		return "", errRegex
	}

	return regex.ReplaceAllString(valor, ""), nil
}

func ImprimirDiferencaEntreDuasDatas(dataHora1 time.Time, dataHora2 time.Time) string {
	hs := dataHora1.Sub(dataHora2).Hours()

	hs, mf := math.Modf(hs)
	ms := mf * 60

	ms, sf := math.Modf(ms)
	ss := sf * 60

	return fmt.Sprint(math.Abs(hs), " horas, ", math.Abs(ms), " minutos, ", math.Abs(ss), " segundos ")
}

func RemoverArquivo(localArquivo string) error {
	if _, err := os.Stat(localArquivo); !os.IsNotExist(err) {
		err := os.Remove(localArquivo)
		if err != nil {
			return fmt.Errorf("Ocorreu um erro ao apagar o arquivo '%s'. Erro: %s", localArquivo, err.Error())
		}
	}

	return nil
}

func RetornaTamanhoArquivo(localArquivo string) (int64, error) {
	fi, err := os.Stat(localArquivo)
	if err != nil {
		return -1, err
	}

	size := fi.Size()

	return size, nil
}

func RetornaValorFormatado(valorEmBytes float64) string {
	var suffixes [5]string

	size := valorEmBytes
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]

	return fmt.Sprint(strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix))
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	newVal = round / pow

	return
}

func RetornaLista(diretorio string) ([]string, error) {
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

func FormataDataPadraoISO8601(DataHora string) string {
	var dataFormatada string

	dataFormatada = fmt.Sprintf("%s%s", strings.Replace(DataHora, " ", "T", 1), ".000-0200")

	return dataFormatada
}

func RemoverItemSlice(Slice interface{}, Indice int) {
	vField := reflect.ValueOf(Slice)

	value := vField.Elem()

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		result := reflect.AppendSlice(value.Slice(0, Indice), value.Slice(Indice+1, value.Len()))

		value.Set(result)
	}
}
