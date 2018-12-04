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
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
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
	URI   string
	resty *resty.Client
}

func NewRestClient(uri string, insecure bool) *RestClient {

	client := http.DefaultClient
	if insecure {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}

		client.Transport = transCfg
	}

	r := resty.NewWithClient(client)
	r.HostURL = uri

	return &RestClient{
		uri,
		r,
	}
}

func (restClient *RestClient) AddAuthHeader(token string) *RestClient {
	/*if headers["Authorization"] == "" {
		//restClient.resty.SetAuthToken(token)
		restClient.resty.SetHeader("Authorization", "bearer "+token)
		restClient.resty.SetHeader("x-anypoint-session-extend", "true")

		headers["Authorization"] = token
	}
	*/
	restClient.resty.SetAuthToken(token)
	return restClient
}

func (restClient *RestClient) AddOrgHeader(orgId string) *RestClient {
	if headers["X-ANYPNT-ORG-ID"] == "" {
		restClient.resty.SetHeader("X-ANYPNT-ORG-ID", orgId)
		headers["X-ANYPNT-ORG-ID"] = orgId
	}
	return restClient
}

func (restClient *RestClient) AddEnvHeader(envId string) *RestClient {
	if headers["X-ANYPNT-ENV-ID"] == "" {
		restClient.resty.SetHeader("X-ANYPNT-ENV-ID", envId)
		headers["X-ANYPNT-ENV-ID"] = envId
	}
	return restClient
}

func (restClient *RestClient) AddHeader(key, value string) *RestClient {
	restClient.resty.SetHeader(key, value)
	return restClient
}

//params should be struct that will be encoded into URL parametrers. Example fron https://godoc.org/github.com/google/go-querystring/query:
//
// type Options struct {
//	Query   string `url:"q"`
//	ShowAll bool   `url:"all"`
//	Page    int    `url:"page"`
// }
//
// opt := Options{ "foo", true, 2 }
// v, _ := query.Values(opt)
// fmt.Print(v.Encode()) // will output: "q=foo&all=true&page=2"
func (restClient *RestClient) GETWithParams(path string, params map[string]string, responseObj interface{}) error {

	Debug(func() {
		log.Println("REQUEST")
		log.Printf("GET %s", restClient.URI+path)
	})

	res, err := restClient.resty.R().SetQueryParams(params).SetResult(&responseObj).Get(path)

	if err != nil {
		fmt.Printf("\nError while performing a GET %s : %s\n", path, err)
		return err
	}

	httpErr := validateResponse(res.RawResponse, err, "GET", path)
	if httpErr != nil {
		Debug(func() {
			fmt.Printf("\nError while performing GET to %q\nError: %s", path, httpErr)
		})
		return httpErr
	}

	if err != nil {
		return fmt.Errorf("Error while reading response for %s : %s ", path, err)
	}

	logResponse("GET", res.Body())

	return nil

}

//GET a resource (no parameters specified in the URI) and fill teh responseObj with a marshalled
//JSON Object
func (restClient *RestClient) GET(path string, responseObj interface{}) error {
	return restClient.GETWithParams(path, nil, responseObj)
}

//PATCH - Perform an HTTP PATCH
func (restClient *RestClient) PATCH(body interface{}, path string, cType ContentType, responseObj interface{}) error {

	Debug(func() {
		log.Println("REQUEST")
		log.Printf("PATCH %s%s", restClient.URI, path)
	})
	res, err := restClient.resty.R().SetBody(body).SetResult(responseObj).Patch(path)

	logResponse("PATCH", res.Body())

	httpErr := validateResponse(res.RawResponse, err, "POST", path)

	return httpErr
}

//POST - Perform an HTTP POST
func (restClient *RestClient) POST(body interface{}, path string, responseObj interface{}) error {
	log.Printf("POST-ing to %s", restClient.URI+path)

	res, err := restClient.resty.R().SetBody(body).SetResult(responseObj).Post(path)

	logRequest(*res.Request)

	if err != nil {
		log.Printf("Error while executing POST %s : %s", path, err)
		return err
	}
	logResponse("POST", res.Body())

	httpErr := validateResponse(res.RawResponse, err, "POST", path)

	return httpErr
}

//PUT - Performs an HTTP PUT
func (restClient *RestClient) PUT(body interface{}, path string, responseObj interface{}) error {
	log.Printf("PUT-ing to %s", restClient.URI+path)

	res, err := restClient.resty.R().SetBody(body).SetResult(responseObj).Put(path)

	logRequest(*res.Request)

	if err != nil {
		log.Printf("Error while executing PUT %s : %s", path, err)
		return err
	}
	logResponse("Put", res.Body())

	httpErr := validateResponse(res.RawResponse, err, "PUT", path)

	return httpErr
}

//DELETE - Perform an HTTP DELETE
func (restClient *RestClient) DELETE(body interface{}, path string, cType ContentType, responseObj interface{}) error {

	Debug(func() {
		log.Println("REQEST")
		log.Printf("DELETE %s%s", restClient.URI, path)
	})

	res, err := restClient.resty.R().SetBody(body).SetResult(responseObj).Delete(path)

	logResponse("DELETE", res.Body())

	httpErr := validateResponse(res.RawResponse, err, "DELETE", path)

	return httpErr
}

func validateResponse(response *http.Response, err error, method, path string) error {

	if err != nil {
		return err
	}

	if response.StatusCode == 401 {
		return NewHttpError(401, "Missing auth token or auth token expired. Please login again.")
	}

	if response.StatusCode == 404 {
		return NewHttpError(404, fmt.Sprintf("Entity %q not found", path))
	}

	if response.StatusCode >= 400 {
		return NewHttpError(response.StatusCode, fmt.Sprintf("\nError when invoking endpoint %s - %s \nHeaders; %s", path, response.Status, response.Request.Header))
	}

	return nil
}

func logRequest(req resty.Request) {
	log.Print("REQUEST DUMP")
	dump, err := httputil.DumpRequest(req.RawRequest, true)
	if err != nil {
		log.Printf("Error while dumping the request.. %s", err)
	}

	if strings.Contains(req.Header.Get("Content-Type"), "application/json") ||
		strings.Contains(req.Header.Get("Content-Type"), "text/json") {
		reqBody, err := json.Marshal(req.Body)
		if err != nil {
			log.Printf("Error while marshalling request body: %s", err)
		} else {
			log.Print(string(reqBody[:]))
		}
	}
	log.Print(string(dump[:]))

	if req.Body != nil {
		log.Print(req.Body)
	}
}

func logResponse(method string, response []byte) {
	log.Printf("RESPONSE")
	log.Printf("%s", response)
}
