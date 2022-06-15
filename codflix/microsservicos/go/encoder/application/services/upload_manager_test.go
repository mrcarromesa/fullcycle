package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func prepare1() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "example.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}

func TestVideoServiceUpload(t *testing.T) {

	video, repo := prepare1()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("codeeducationcodflixtest")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "codeeducationcodflixtest"
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(5, doneUpload)

	result := <-doneUpload

	require.Equal(t, "uploaded complete", result)

}
