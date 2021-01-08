package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bdavs3/worker/server/api"
	"github.com/bdavs3/worker/worker"
)

// TODO (out of scope): Rather than using hard-coded user credentials, provide
// the user with a way to create an account and log in. Once authenticated
// with the API, the client could receive a session token that precludes the
// need to authenticate on each subsequent request.

const (
	crtFile = "../worker.crt"
	host    = "https://localhost:443"
	timeout = 5 * time.Second
)

// Client is an HTTP client for making requests to control and assess Linux
// processes on the server. Use NewClient to create a new instance.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Client instance.
func NewClient() (*Client, error) {
	rootCAs, err := generateCertPool(crtFile)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL: host,
		HTTPClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: rootCAs,
				},
			},
		},
	}

	return client, nil
}

// generaeteCertPool returns a CertPool containing the given certificate.
func generateCertPool(crtFile string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(crtFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return caCertPool, err
}

// PostJob passes a Linux process to the worker library for execution.
func (c *Client) PostJob(job worker.Job) (string, error) {
	requestBody, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	response, err := c.makeRequestWithAuth(
		http.MethodPost,
		"/jobs/run",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}

	return response.ID, nil
}

// GetJobStatus queries the status of a process being handled by the worker library.
func (c *Client) GetJobStatus(id string) (string, error) {
	response, err := c.makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("/jobs/%s/status", id),
		nil,
	)
	if err != nil {
		return "", err
	}
	if response.ID != id {
		return "", errors.New("response contains incorrect job id")
	}

	return response.Status, nil
}

// GetJobOutput queries the output of a process being handled by the worker library.
func (c *Client) GetJobOutput(id string) (string, error) {
	response, err := c.makeRequestWithAuth(
		http.MethodGet,
		fmt.Sprintf("/jobs/%s/out", id),
		nil,
	)
	if err != nil {
		return "", err
	}
	if response.ID != id {
		return "", errors.New("response contains incorrect job id")
	}

	return response.Output, nil
}

// KillJob terminates a process being handled by the worker library.
func (c *Client) KillJob(id string) (string, error) {
	response, err := c.makeRequestWithAuth(
		http.MethodPut,
		fmt.Sprintf("/jobs/%s/kill", id),
		nil,
	)
	if err != nil {
		return "", err
	}
	if response.ID != id {
		return "", errors.New("response contains incorrect job id")
	}

	return response.Status, nil
}

// makeRequestWithAuth makes an HTTP request to the given endpoint
// after setting the Authorization header. It then returns the response.
func (c *Client) makeRequestWithAuth(method, endpoint string, requestBody io.Reader) (*api.Response, error) {
	req, err := http.NewRequest(
		method,
		c.BaseURL+endpoint,
		requestBody,
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("username"), os.Getenv("pw"))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s\n%s", http.StatusText(resp.StatusCode), body)
	}

	var response *api.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
