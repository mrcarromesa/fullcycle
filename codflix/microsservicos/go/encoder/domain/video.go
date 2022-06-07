package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResoruceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path" valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"-" valid:"-"`
	Jobs       []*Job    `json:"-" valid:"-" gorm:"ForeignKey:VideoID"` // ForeignKey o campo VideoID precisa existir na estrutura do job.go Podemos ter vários jobs dentro de vídeo
}

// A função init roda antes do que qualquer coisa
func init() {
	// Garante que tudo que é obrigatório ser preenchido por padrão
	validator.SetFieldsRequiredByDefault(true)
}

// Estamos passando por referencia, e todas as vezes que alguém gerar o NewVideo e alterar essa variavel, ela será alterada em qualquer lugar.
func NewVideo() *Video {
	return &Video{}
}

// Escrevendo dessa forma diferente fazemos que Validate seja um metódo da struct Video.
func (video *Video) Validate() error {

	// valid, error := validator.ValidateStruct(video)
	// Estamos utilizando um blank identifier o _
	_, error := validator.ValidateStruct(video)

	if error != nil {
		return error
	}

	return nil
}
