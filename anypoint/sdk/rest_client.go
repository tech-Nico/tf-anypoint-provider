// Copyright © 2017 Nico Balestra <functions@protonmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"crypto/tls"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/tech-nico/anypoint-cli/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

type HttpError struct {
	StatusCode int
	msg        string
}

type ContentType int

var headers map[string]string = make(map[string]string, 0)

const (
	Application_Json ContentType = iota
	Application_OctetStream
	Application_Pdf
	Application_Atom_Xml
	Application_Form_Urlencoded
	Application_SVG_XML
	Application_XHTML_XML
	Application_XML
	Multipart_Form_Data
	Text_HTML
	Text_Plain
	Text_XML
	Wildcard
)

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP Error %d - %s", e.StatusCode, e.msg)
}

func NewHttpError(code int, theMsg string) error {
	return &HttpError{
		StatusCode: code,
		msg:        theMsg,
	}
}

type RestClient struct {
	URI    string
	Sling  *sling.Sling
	client *http.Client
}

func NewRestClient(uri string, insecure bool) *RestClient {

	client := &http.Client{}
	if insecure {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}

		client.Transport = transCfg
	}

	s := sling.New().
		Client(client).
		Base(uri)

	return &RestClient{
		uri,
		s,
		client,
	}
}

func (client *RestClient) AddAuthHeader(token string) *RestClient {
	if headers["Authorization"] == "" {
		client.Sling.Add("Authorization", "Bearer "+token)
		headers["Authorization"] = token
	}
	return client
}

func (client *RestClient) AddOrgHeader(orgId string) *RestClient {
	if headers["X-ANYPNT-ORG-ID"] == "" {
		client.Sling.Add("X-ANYPNT-ORG-ID", orgId)
		headers["X-ANYPNT-ORG-ID"] = orgId
	}
	return client
}

func (client *RestClient) AddEnvHeader(envId string) *RestClient {
	if headers["X-ANYPNT-ENV-ID"] == "" {
		client.Sling.Add("X-ANYPNT-ENV-ID", envId)
		headers["X-ANYPNT-ENV-ID"] = envId
	}
	return client
}

func (client *RestClient) AddHeader(key, value string) *RestClient {
	client.Sling.Add(key, value)
	return client
}

func (client *RestClient) GET(path string) ([]byte, error) {

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("GET %s", client.URI+path)
	})
	req, err := client.Sling.Get(path).Request()
	if err != nil {
		fmt.Printf("\nError building GET request for path %s : %s\n", path, err)
		return nil, err
	}
	res, err := client.client.Do(req)
	defer res.Body.Close()

	httpErr := validateResponse(res, err, "GET", path)
	if httpErr != nil {
		utils.Debug(func() {
			fmt.Printf("\nError while performing GET to %q\nError: %s", path, httpErr)
		})
		return nil, httpErr
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response for %s : %s ", path, err)
	}

	utils.Debug(logResponse("GET", res))

	return body, nil
}

//PATCH - Perform an HTTP PATCH
func (client *RestClient) PATCH(body interface{}, path string, cType ContentType, responseObj interface{}) (*http.Response, error) {

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("PATCH %s%s", client.URI, path)
	})
	sling := client.Sling.Patch(path)
	sling = setSlingBodyForContentType(cType, sling, body)

	response, err := sling.ReceiveSuccess(responseObj)

	utils.Debug(logResponse("PATCH", response))

	httpErr := validateResponse(response, err, "POST", path)

	return response, httpErr
}

//POST - Perform an HTTP POST
func (client *RestClient) POST(body interface{}, path string, cType ContentType, responseObj interface{}) (*http.Response, error) {

	sling := client.Sling.Post(path)

	sling = setSlingBodyForContentType(cType, sling, body)
	req, err := sling.Request()
	utils.Debug(logRequest(req))

	response, err := sling.ReceiveSuccess(responseObj)

	utils.Debug(logResponse("POST", response))

	httpErr := validateResponse(response, err, "POST", path)

	return response, httpErr
}

//DELETE - Perform an HTTP DELETE
func (client *RestClient) DELETE(body interface{}, path string, cType ContentType, responseObj interface{}) (*http.Response, error) {

	utils.Debug(func() {
		log.Println("REQEST")
		log.Printf("POST %s%s", client.URI, path)
	})

	sling := client.Sling.Delete(path)

	sling = setSlingBodyForContentType(cType, sling, body)

	response, err := sling.ReceiveSuccess(responseObj)

	utils.Debug(logResponse("DELETE", response))

	httpErr := validateResponse(response, err, "DELETE", path)

	return response, httpErr
}

func setSlingBodyForContentType(cType ContentType, sling *sling.Sling, body interface{}) *sling.Sling {
	if body != nil {
		switch cType {
		case Application_Json:
			sling = sling.BodyJSON(body)
		case Application_Form_Urlencoded:
			sling = sling.BodyForm(body)
		case Application_OctetStream:
			sling = sling.Body(body.(io.Reader))
		default:
			sling = sling.Body(body.(io.Reader))
		}
	}
	return sling
}

func validateResponse(response *http.Response, err error, method, path string) error {

	if err != nil {
		return err
	}

	if response.StatusCode == 401 {
		return NewHttpError(401, "Auth token expired. Please login again")
	}

	if response.StatusCode == 404 {
		return NewHttpError(404, fmt.Sprintf("Entity %q not found", path))
	}

	if response.StatusCode >= 400 {
		return NewHttpError(response.StatusCode, fmt.Sprintf("\nError when invoking endpoint %s - %s \nHeaders; %s", path, response.Status, response.Request.Header))
	}

	return nil
}

func logRequest(req *http.Request) func() {
	return func() {
		log.Print("REQUEST DUMP")
		dump, _ := httputil.DumpRequest(req, true)
		log.Print(dump)
	}
}

func logResponse(method string, response *http.Response) func() {
	return func() {
		log.Printf("RESPONSE")
		dump, _ := httputil.DumpResponse(response, true)
		log.Printf("\n %s RESPONSE: %", method, dump)
	}

}
