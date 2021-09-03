package attachas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/writeas/impart"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type (
	// File represents a photo on Snap.as.
	File struct {
		ID       string    `json:"id"`
		Created  time.Time `json:"created"`
		Filename string    `json:"filename"`
		Size     int64     `json:"size"`
		URL      string    `json:"url"`
	}

	// FileParams holds valid values for uploading files.
	FileParams struct {
		FileName string
		Body     string

		// OrgAlias is the alias of the organization to upload to.
		OrgAlias string
	}
)

// UploadFile uploads a photo, and returns an Attach.as File.
func (c *Client) UploadFile(sp *FileParams) (*File, error) {
	f, err := os.Open(sp.FileName)
	if err != nil {
		return nil, fmt.Errorf("open file: %s", err)
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	if sp.Body != "" {
		err = w.WriteField("body", sp.Body)
		if err != nil {
			return nil, fmt.Errorf("write field 'body': %s", err)
		}
	}

	part, err := w.CreateFormFile("file", filepath.Base(f.Name()))
	if err != nil {
		return nil, fmt.Errorf("create form file: %s", err)
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return nil, fmt.Errorf("copy file: %s", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer: %s", err)
	}

	orgBase := ""
	if sp.OrgAlias != "" {
		orgBase = "/organizations/" + sp.OrgAlias
	}
	url := fmt.Sprintf("%s%s%s", c.Config.BaseURL, orgBase, "/files/upload")
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %s", err)
	}
	req.Header.Add("User-Agent", c.Config.UserAgent)
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("Authorization", c.Token)

	resp, err := c.Config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %s", err)
	}
	defer resp.Body.Close()

	env := &impart.Envelope{
		Code: resp.StatusCode,
		Data: &File{},
	}
	err = json.NewDecoder(resp.Body).Decode(&env)
	if err != nil {
		return nil, err
	}
	if env.Code != http.StatusCreated {
		return nil, fmt.Errorf("%s", env.ErrorMessage)
	}

	return env.Data.(*File), nil
}
