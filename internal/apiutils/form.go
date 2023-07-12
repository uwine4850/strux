package apiutils

import (
	"bytes"
	"encoding/json"
	"github.com/uwine4850/strux_api/services/protofiles/baseproto"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// NewApiForm The structure that sends requests to the api.
// Each key in the TextValues and FileValues maps stands for the name of the field in the form.
// TextValues is responsible for the text values of the form.
// FileValues is for file values, namely the slice contains the full absolute paths to the files.
type NewApiForm struct {
	Method     string
	Url        string
	TextValues map[string]string
	FileValues map[string][]string
}

// SendForm method for sending the form
func (f *NewApiForm) SendForm() (*baseproto.BaseResponse, *multipart.Form, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := f.writeTextField(writer)
	if err != nil {
		return nil, nil, err
	}
	err = f.writeFiles(writer)
	if err != nil {
		return nil, nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}

	r, err := http.NewRequest(f.Method, f.Url, body)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, nil, err
	}
	res, form, err := f.getResponse(resp)
	if err != nil {
		return nil, nil, err
	}
	return res, form, nil
}

// getResponse Getting a response from the server after submitting the form.
// The response can be of two kinds:
// Response in json format.
// Response in form format.
func (f *NewApiForm) getResponse(resp *http.Response) (*baseproto.BaseResponse, *multipart.Form, error) {
	contentType := resp.Header.Get("Content-Type")
	mt, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, nil, err
	}
	// get base response(json)
	if mt == "text/plain" {
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		baseResponse := &baseproto.BaseResponse{}
		err = json.Unmarshal(all, baseResponse)
		if err != nil {
			return nil, nil, err
		}
		return baseResponse, nil, nil
	}
	boundary, _ := params["boundary"]
	reader := multipart.NewReader(resp.Body, boundary)
	form, err := reader.ReadForm(64 << 20)
	if err != nil {
		return nil, nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	return nil, form, err
}

// writeTextField writes text fields to the form
func (f *NewApiForm) writeTextField(writer *multipart.Writer) error {
	for fieldName, fieldValue := range f.TextValues {
		err := writer.WriteField(fieldName, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeFiles writes file fields to the form. One field can contain multiple files.
func (f *NewApiForm) writeFiles(writer *multipart.Writer) error {
	for fieldName, files := range f.FileValues {
		for i := 0; i < len(files); i++ {
			file, err := os.Open(files[i])
			if err != nil {
				return err
			}
			formFile, err := writer.CreateFormFile(fieldName, filepath.Base(files[i]))
			if err != nil {
				return err
			}
			_, err = io.Copy(formFile, file)
			if err != nil {
				return err
			}
			err = file.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
