# Microsserviço para encoder de vídeo

- Dockerfile para o projeto obtido de: https://gist.github.com/wesleywillians/9dcee3aa242ffb6bc92e7f0fdbc7aadd

- Para iniciar rodar o comando:

```shell
docker compose up --build
```

- Para acessar o container executar o comando:

```shell
docker compose exec app bash
```

---

## go mod

- utilizamos para gerenciar as versões dos nosso projeto

- Criamos a pasta `encoder/` e dentro dela rodar o comando:

```shell
go mod init encoder
```

- Com isso será criado um arquivo `go.mod`

- E o go.mod serve para gerenciar as nossas dependencias.

- Ao executar o comando:

```shell
go run main.go
```

- Caso o pacote não esteja instalado ele irá baixar a dependencia para nós automáticamente ou seja ele irá realizar um `go get` do pacote

- E ele gera um arquivo `go.sum` que é um arquivo lock das dependencias


---

### Testes

- Pacote para testes: github.com/stretchr/testify/require

- Para realizar um test primeiro criamos o arquivo `encoder/domain/video_test.go`
- Para executar os testes utilizar o comando dentro da pasta  raiz:

```shell
go test ./...
```

- Pois assim ele irá executar todos os testes que encontrar.

- Muitas vezes é interesante limpar o cache dos testes podemos fazer isso executando esse comando:

```shell
go clean -testcache
```

---

### Validation

- Pacote para validação: github.com/asaskevich/govalidator

- Exemplo de uso:

```go
type Job struct {
	ID               string    `valid:"uuid"`
	OutputBucketPath string    `valid:"notnull"`
	Status           string    `valid:"notnull"`
	Video            *Video    `valid:"-"` // O * é um ponteiro para um objeto no caso o objeto Video ou seja aponta para o mesmo local da memória do objeto Video
	VideoID          string    `valid:"-"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `valid:"-"`
	UpdatedAt        time.Time `valid:"-"`
}
```

---

### UUID

- Pacote para gerar uuid: github.com/satori/go.uuid


---

### Visibilidade

- Quando eu quero que uma func seja visivel fora do pacote inicio ela com letra maiuscula, caso queira que ela seja visivel fora do package utilizo a primeira letra maiuscula


---

### ORM

- Lib para utilizar: github.com/jinzhu/gorm

- Exemplo de conexão: `encoder/framework/database/db.go` e `encoder/.env`

- Utiliziamos a anotação gorm:

- Definindo os campos:

```go
type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResoruceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path" valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"-" valid:"-"`
	Jobs       []*Job    `json:"-" valid:"-" gorm:"ForeignKey:VideoID"` // ForeignKey Podemos ter vários jobs dentro de vídeo
}
```

- Quando adicionamos uma relação precisamos adicionar no `framework/database/db.go > Connect > AutoMigrateDB`:

```go
if d.AutoMigrateDB {
		d.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.DB.Model(domain.Job{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "CASCADE")
	}
```

- Dentro de `framework/database/db.go` precisamos utilizar algumas importações, e o go não permite importar coisas que não estão sendo utilizadas no código então para isso utilizamos o `_` na frente da importação dessa forma:

```go
import (
	// ... MORE
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pg"
)
```

---

### Converter nome de campos quando for JSON

- Utilizamos a anotação `json:"NOME_DO_CAMPO"`

```go
type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Video            *Video    `json:"video" valid:"-"`
	VideoID          string    `json:"-"` // NÃO PRECISAMOS DELE POIS ELE ESTARÁ INSERIDO NO CAMPO `video: {}`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}
```

---

### Repositories

- Utilizar para lidar com as chamadas ao banco de dados

- para tal criamos a pasta `encoder/application/repositories`


---

### Baixando dependencias go

- Quando precisamos baixar todas as dependencias utilizando go, podemos remover o arquivo `go.sum` e executar o comando:

```shell
go mod download
```

---

### Criar um account service no google

- Acessar o IAM > Contas de serviço
- Cliar para criar nova conta de serviço
- dar um nome para ela no passo 1
- No passo 2 Adicionar o papel de `Administador de ambiente e de objetos do Storage`
- No passo 3 manter como está e concluir.
- Editar a conta e serviço e ir na aba Chaves e criar uma nova chave no formato JSON e realizar o download da mesma para o projeto. 

- Adicionar no arquiv `.env` o seguinte:

```
GOOGLE_APPLICATION_CREDENTIALS="ARQUIVO_CRENDENTIAL.json"
```

- No lugar de `ARQUIVO_CRENDENTIAL.json` colocar o nome do seu arquivo

---

- Para utilizar as funcionalidades do gcp utilizamos a dependencia: `cloud.google.com/go/storage`

- Para realizar upload é necessario alterar o ACL de uniforme para detalhado

---

### Upload

- Para realização do upload será utilizado o `encoder/application/services/upload_manager.go`
- Ele irá obter os Paths dos arquivos
- Irá conectar-se com o GCP
- Irá utilizar o go routines para realizar o upload de todos os paths ou arquivos.
- Ao terminar o upload ou obter erro ele encerra o canal e go routine

---

### Jobs

- Será através dos jobs que serão feitas as etapas de:
- Baixar o arquivo do GCP
- Fragmentar o arquivo
- Realizar o encode
- Por fim realizar o upload do arquivo

- Para cada uma dessas etapas será salvo/atualizado o status do processo na base de dados, 
- Também será tratado o erro em cada uma das etapas caso venha a ocorrer

- Esse job será chamado via fila para gerenciar essas filas será utilizado o RabbitMQ

- O RabbitMQ combina mais com a camada de framework