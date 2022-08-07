package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": "234",
		"file_path": "file.mp4",
		"status": "pending"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `abc`
	err = utils.IsJson(json)
	require.Error(t, err)

}
