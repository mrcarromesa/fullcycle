package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {

	// caminho/x/b/arquivo.mp4
	// split:
	// [0] caminho/x/b/
	// [1] arquivo.mp4

	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)

	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

// Verificar se os arquivos existem e adicionar o caminho de cada um deles no array Paths
func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// Adiciona o caminho dos arquivos no array Paths
			vu.Paths = append(vu.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {

	// Canal que informa qual item ele deve ler.
	in := make(chan int, runtime.NumCPU()) // qual o arquivo baseado na posicao do slice Paths
	returnChannel := make(chan string)     // canal que fica monitorando o retorno

	err := vu.loadPaths() // carrega os paths para VideoUpload.Paths
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload() // pega o cliente do google storage

	if err != nil {
		return err
	}

	// o concurrency é a quantidade de trheads em paralelo que queremos que ocorra
	for process := 0; process < concurrency; process++ {
		// Adicionamos o go na frente para que ela seja executada em uma go routine
		// Uma vez iniciada a go routina ela ficará rodando para sempre, digamos assim. até que o canal seja fechado
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx) // Será iniciada novas rotinas
		// Esse routine irá ler o in que é a posição dos paths, uma vez que ele ler esvazia novamente.
		// Ele tem que estar vazio para que seja feita uma nova entrada e leitura
	}

	// funcao go anonima
	go func() {
		// ira pegar a posicao arquivo por arquivo e setar no in
		for x := 0; x < len(vu.Paths); x++ {
			// Estou passando o valor de x para variavel in
			in <- x
			// o in terá o valor do path no caso a posicao dele, e só poderá receber um novo valor quando ele for esvazido ou seja x ficará aguardando até que o in seja esvaziado
			// a função vu.uploadworker irá esvaziar o valor de in `Ao ler o valor por exemplo`, liberando ele para receber um novo valor de x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			// outro canal para informar que o upload deu errado
			// E temos que para tudo pois nao adianta um arquivo ter dado certo outro errado, entao para tudo.
			// O quando tudo foi concluído com sucesso daí pedimos para finalizar também o canal!
			doneUpload <- r
			break
		}
	}

	return nil

}

// Funcao de upload
// terá varias threads para fazer uploads e essa funcao será as threads!
// a utilizacao do chan = canal + int ou seja é um canal que aceita/recebe inteiro
func (vu *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during the upload: %v. Error: %v", vu.Paths[x], err)

			// retornamos para o canal o erro
			returnChan <- err.Error()
		}

		// Se receber uma string em branco significa que não houve erro -- vamos tratar assim
		returnChan <- ""

	}

	// utilizamos isso para fechar o canal. e finalizar o doneUpload
	returnChan <- "uploaded complete"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
