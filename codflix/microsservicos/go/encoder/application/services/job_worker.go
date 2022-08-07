package services

import (
	"encoder/domain"
	"encoder/framework/utils"
	"encoding/json"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job domain.Job, workerID int) {

	/**
	// receber o json:
	{
		"resource_id": "id do video da pessoa que enviou para nossa fila",
		"file_path": "caminho do arquivo no bucket"
	}
	*/

	for message := range messageChannel {
		// pegar message body do json.
		// validar o json é valido?
		// validar video
		// inserir video no banco de dados
		// start

		err := utils.IsJson(string(message.Body))

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		jobService.VideoService.Video.ID = uuid.NewV4().String()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.Video.Validate()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("OUTPUT_BUCKET_NAME")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		_, err = jobService.JobRepository.Insert(&job)

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job
		err = jobService.Start()

		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(job, message, nil)

	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}

	return result
}
