package vagrantcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/packer/common"
)

type VagrantCloudClient struct {
	// The http client for communicating
	client *http.Client

	// The base URL of the API
	BaseURL string

	// Access token
	AccessToken string
}

type VagrantCloudErrors struct {
	Errors map[string][]string `json:"errors"`
}

func (v VagrantCloudErrors) FormatErrors() string {
	errs := make([]string, 0)
	for e := range v.Errors {
		msg := fmt.Sprintf("%s %s", e, strings.Join(v.Errors[e], ","))
		errs = append(errs, msg)
	}
	return strings.Join(errs, ". ")
}

func (v VagrantCloudClient) New(baseUrl string, token string) *VagrantCloudClient {
	c := &VagrantCloudClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
		BaseURL:     baseUrl,
		AccessToken: token,
	}
	return c
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

func (v VagrantCloudClient) Get(path string) (*http.Response, error) {
	params := url.Values{}
	params.Set("access_token", v.AccessToken)
	reqUrl := fmt.Sprintf("%s/%s?%s", v.BaseURL, path, params.Encode())

	// Scrub API key for logs
	scrubbedUrl := strings.Replace(reqUrl, v.AccessToken, "ACCESS_TOKEN", -1)
	log.Printf("Post-Processor Vagrant Cloud API GET: %s", scrubbedUrl)

	req, err := http.NewRequest("GET", reqUrl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := v.client.Do(req)

	log.Printf("Post-Processor Vagrant Cloud API Response: \n\n%+v", resp)

	return resp, err
}

func (v VagrantCloudClient) Delete(path string) (*http.Response, error) {
	params := url.Values{}
	params.Set("access_token", v.AccessToken)
	reqUrl := fmt.Sprintf("%s/%s?%s", v.BaseURL, path, params.Encode())

	// Scrub API key for logs
	scrubbedUrl := strings.Replace(reqUrl, v.AccessToken, "ACCESS_TOKEN", -1)
	log.Printf("Post-Processor Vagrant Cloud API DELETE: %s", scrubbedUrl)

	req, err := http.NewRequest("DELETE", reqUrl, nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := v.client.Do(req)

	log.Printf("Post-Processor Vagrant Cloud API Response: \n\n%+v", resp)

	return resp, err
}

func (v VagrantCloudClient) Upload(path string, url string, output func(string)) (*http.Response, error) {
	// Open up the file for upload
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file for upload: %s", err)
	}

	// Grab the file size and update the progress bar
	fi, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Error stating file for upload: %s", err)
	}

	// Grab the default looking ProgressBar
	pb := common.GetDefaultProgressBar() // from common/step_download.go
	pb.Total = fi.Size()
	log.Printf("Post-Processor Vagrant Cloud Status: Creating progress-bar and setting total size to %d.\n", pb.Total)

	// Prepare it and set it's output callback
	bar := pb.Start()
	defer bar.Finish()
	bar.Callback = output
	log.Printf("Post-Processor Vagrant Cloud Status: Started progress-bar with output set to %#v.\n", bar.Callback)

	// Prepare the http request with a ProxyReader for the ProgressBar
	proxyFileReader := bar.NewProxyReader(file)
	request, err := http.NewRequest("PUT", url, proxyFileReader)
	if err != nil {
		return nil, fmt.Errorf("Error preparing upload request: %s", err)
	}
	request.ContentLength = fi.Size()

	// Now we can upload the file
	log.Printf("Post-Processor Vagrant Cloud API Upload: %s %s", path, url)
	output("Making Post-Processor Vagrant Cloud upload request. Progress bar should look like : %s", bar.String())
	resp, err := v.client.Do(request)
	output("Completed Post-Processor Vagrant Cloud upload request.")

	// Log the response and we're done.
	log.Printf("Post-Processor Vagrant Cloud Upload Response: \n\n%+v", resp)
	return resp, err
}

func (v VagrantCloudClient) Post(path string, body interface{}) (*http.Response, error) {
	params := url.Values{}
	params.Set("access_token", v.AccessToken)
	reqUrl := fmt.Sprintf("%s/%s?%s", v.BaseURL, path, params.Encode())

	encBody, err := encodeBody(body)

	if err != nil {
		return nil, fmt.Errorf("Error encoding body for request: %s", err)
	}

	// Scrub API key for logs
	scrubbedUrl := strings.Replace(reqUrl, v.AccessToken, "ACCESS_TOKEN", -1)
	log.Printf("Post-Processor Vagrant Cloud API POST: %s. \n\n Body: %s", scrubbedUrl, encBody)

	req, err := http.NewRequest("POST", reqUrl, encBody)
	req.Header.Add("Content-Type", "application/json")

	resp, err := v.client.Do(req)

	log.Printf("Post-Processor Vagrant Cloud API Response: \n\n%+v", resp)

	return resp, err
}

func (v VagrantCloudClient) Put(path string) (*http.Response, error) {
	params := url.Values{}
	params.Set("access_token", v.AccessToken)
	reqUrl := fmt.Sprintf("%s/%s?%s", v.BaseURL, path, params.Encode())

	// Scrub API key for logs
	scrubbedUrl := strings.Replace(reqUrl, v.AccessToken, "ACCESS_TOKEN", -1)
	log.Printf("Post-Processor Vagrant Cloud API PUT: %s", scrubbedUrl)

	req, err := http.NewRequest("PUT", reqUrl, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := v.client.Do(req)

	log.Printf("Post-Processor Vagrant Cloud API Response: \n\n%+v", resp)

	return resp, err
}
