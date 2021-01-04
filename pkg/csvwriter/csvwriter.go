package csvwriter

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFile(w http.ResponseWriter, r *http.Request) (string, error) {

	r.ParseMultipartForm(10 << 20) // max 10mb files

	// retrieve file from posted form-data
	file, _, err := r.FormFile("user_csv")
	if err != nil {
		return "", fmt.Errorf("Error retrieving file from form-data (%v)", err)
	}
	defer file.Close()

	//  write temporary file
	tempFile, err := ioutil.TempFile("tempcsv", "upload-*.csv")
	if err != nil {
		return "", fmt.Errorf("Error writing temp file (%v)", err)
	}
	defer tempFile.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error copying data to tempfile (%v)", err)
	}
	tempFile.Write(fileContent)
	filepath := tempFile.Name()

	return filepath, nil
}
