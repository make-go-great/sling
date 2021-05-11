package sling

import (
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	slinghttp "github.com/haunt98/sling/internal/http"
)

// Sling is an HTTP Request builder and sender.
type Sling struct {
	httpClient      slinghttp.Client
	method          string
	pathURL         *url.URL
	header          http.Header
	queryStructs    []interface{}
	bodyProvider    BodyProvider
	responseDecoder ResponseDecoder
}

// New returns a new Sling with an http DefaultClient.
func New() *Sling {
	return &Sling{
		httpClient:      http.DefaultClient,
		method:          "GET",
		header:          make(http.Header),
		queryStructs:    make([]interface{}, 0),
		responseDecoder: jsonDecoder{},
	}
}

// New returns a copy of a Sling for creating a new Sling with properties
// from a parent Sling. For example,
//
// 	parentSling := sling.New().Client(client).Base("https://api.io/")
// 	fooSling := parentSling.New().Get("foo/")
// 	barSling := parentSling.New().Get("bar/")
//
// fooSling and barSling will both use the same client, but send requests to
// https://api.io/foo/ and https://api.io/bar/ respectively.
//
// Note that query and body values are copied so if pointer values are used,
// mutating the original value will mutate the value within the child Sling.
func (s *Sling) New() *Sling {
	// copy Headers pairs into new Header map
	headerCopy := make(http.Header)
	for k, v := range s.header {
		headerCopy[k] = v
	}
	return &Sling{
		httpClient:      s.httpClient,
		method:          s.method,
		pathURL:         s.pathURL,
		header:          headerCopy,
		queryStructs:    append([]interface{}{}, s.queryStructs...),
		bodyProvider:    s.bodyProvider,
		responseDecoder: s.responseDecoder,
	}
}

// HTTP client

func (s *Sling) HTTPClient(client slinghttp.Client) *Sling {
	if client == nil {
		s.httpClient = http.DefaultClient
		return s
	}

	s.httpClient = client
	return s
}

// HTTP Method
// See: https://golang.org/pkg/net/http/#pkg-constants

func (s *Sling) Get(pathURL string) *Sling {
	return s.Method(http.MethodGet, pathURL)
}

func (s *Sling) Head(pathURL string) *Sling {
	return s.Method(http.MethodHead, pathURL)
}

func (s *Sling) Post(pathURL string) *Sling {
	return s.Method(http.MethodPost, pathURL)
}

func (s *Sling) Put(pathURL string) *Sling {
	return s.Method(http.MethodPut, pathURL)
}

func (s *Sling) Patch(pathURL string) *Sling {
	return s.Method(http.MethodPatch, pathURL)
}

func (s *Sling) Delete(pathURL string) *Sling {
	return s.Method(http.MethodDelete, pathURL)
}

func (s *Sling) Connect(pathURL string) *Sling {
	return s.Method(http.MethodConnect, pathURL)
}

func (s *Sling) Options(pathURL string) *Sling {
	return s.Method(http.MethodOptions, pathURL)
}

func (s *Sling) Trace(pathURL string) *Sling {
	return s.Method(http.MethodTrace, pathURL)
}

func (s *Sling) Method(method, pathURL string) *Sling {
	s.method = method
	return s.PathURL(pathURL)
}

// Header

func (s *Sling) AddHeader(key, value string) *Sling {
	s.header.Add(key, value)
	return s
}

func (s *Sling) SetHeader(key, value string) *Sling {
	s.header.Set(key, value)
	return s
}

// URL

func (s *Sling) PathURL(urlStr string) *Sling {
	pathURL, err := url.Parse(urlStr)
	if err != nil {
		return s
	}

	s.pathURL = pathURL
	return s
}

// QueryStruct appends the queryStruct to the Sling's queryStructs. The value
// pointed to by each queryStruct will be encoded as url query parameters on
// new requests (see Request()).
// The queryStruct argument should be a pointer to a url tagged struct. See
// https://godoc.org/github.com/google/go-querystring/query for details.
func (s *Sling) QueryStruct(queryStruct interface{}) *Sling {
	if queryStruct != nil {
		s.queryStructs = append(s.queryStructs, queryStruct)
	}
	return s
}

// Body

// Body sets the Sling's body. The body value will be set as the Body on new
// requests (see Request()).
// If the provided body is also an io.Closer, the request Body will be closed
// by http.Client methods.
func (s *Sling) Body(body io.Reader) *Sling {
	if body == nil {
		return s
	}
	return s.BodyProvider(bodyProvider{body: body})
}

// BodyProvider sets the Sling's body provider.
func (s *Sling) BodyProvider(body BodyProvider) *Sling {
	if body == nil {
		return s
	}
	s.bodyProvider = body

	ct := body.ContentType()
	if ct != "" {
		s.SetHeader(contentType, ct)
	}

	return s
}

// BodyJSON sets the Sling's bodyJSON. The value pointed to by the bodyJSON
// will be JSON encoded as the Body on new requests (see Request()).
// The bodyJSON argument should be a pointer to a JSON tagged struct. See
// https://golang.org/pkg/encoding/json/#MarshalIndent for details.
func (s *Sling) BodyJSON(bodyJSON interface{}) *Sling {
	if bodyJSON == nil {
		return s
	}
	return s.BodyProvider(jsonBodyProvider{payload: bodyJSON})
}

// BodyForm sets the Sling's bodyForm. The value pointed to by the bodyForm
// will be url encoded as the Body on new requests (see Request()).
// The bodyForm argument should be a pointer to a url tagged struct. See
// https://godoc.org/github.com/google/go-querystring/query for details.
func (s *Sling) BodyForm(bodyForm interface{}) *Sling {
	if bodyForm == nil {
		return s
	}
	return s.BodyProvider(formBodyProvider{payload: bodyForm})
}

// Requests

// Request returns a new http.Request created with the Sling properties.
// Returns any errors parsing the rawURL, encoding query structs, encoding
// the body, or creating the http.Request.
func (s *Sling) Request() (*http.Request, error) {
	err := addQueryStructs(s.pathURL, s.queryStructs)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if s.bodyProvider != nil {
		body, err = s.bodyProvider.Body()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(s.method, s.pathURL.String(), body)
	if err != nil {
		return nil, err
	}
	addHeaders(req, s.header)
	return req, err
}

// addQueryStructs parses url tagged query structs using go-querystring to
// encode them to url.Values and format them onto the url.RawQuery. Any
// query parsing or encoding errors are returned.
func addQueryStructs(reqURL *url.URL, queryStructs []interface{}) error {
	urlValues, err := url.ParseQuery(reqURL.RawQuery)
	if err != nil {
		return err
	}
	// encodes query structs into a url.Values map and merges maps
	for _, queryStruct := range queryStructs {
		queryValues, err := query.Values(queryStruct)
		if err != nil {
			return err
		}
		for key, values := range queryValues {
			for _, value := range values {
				urlValues.Add(key, value)
			}
		}
	}
	// url.Values format to a sorted "url encoded" string, e.g. "key=val&foo=bar"
	reqURL.RawQuery = urlValues.Encode()
	return nil
}

// addHeaders adds the key, value pairs from the given http.Header to the
// request. Values for existing keys are appended to the keys values.
func addHeaders(req *http.Request, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

// Sending

// ResponseDecoder sets the Sling's response decoder.
func (s *Sling) ResponseDecoder(decoder ResponseDecoder) *Sling {
	if decoder == nil {
		return s
	}
	s.responseDecoder = decoder
	return s
}

// ReceiveSuccess creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV.
// Any error creating the request, sending it, or decoding a 2XX response
// is returned.
func (s *Sling) ReceiveSuccess(successV interface{}) (*http.Response, error) {
	return s.Receive(successV, nil)
}

// Receive creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV and
// other responses are JSON decoded into the value pointed to by failureV.
// If the status code of response is 204(no content) or the Content-Lenght is 0,
// decoding is skipped. Any error creating the request, sending it, or decoding
// the response is returned.
// Receive is shorthand for calling Request and Do.
func (s *Sling) Receive(successV, failureV interface{}) (*http.Response, error) {
	req, err := s.Request()
	if err != nil {
		return nil, err
	}
	return s.Do(req, successV, failureV)
}

// Do sends an HTTP request and returns the response. Success responses (2XX)
// are JSON decoded into the value pointed to by successV and other responses
// are JSON decoded into the value pointed to by failureV.
// If the status code of response is 204(no content) or the Content-Length is 0,
// decoding is skipped. Any error sending the request or decoding the response
// is returned.
func (s *Sling) Do(req *http.Request, successV, failureV interface{}) (*http.Response, error) {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	// when err is nil, resp contains a non-nil resp.Body which must be closed
	defer resp.Body.Close()

	// The default HTTP client's Transport may not
	// reuse HTTP/1.x "keep-alive" TCP connections if the Body is
	// not read to completion and closed.
	// https://golang.org/pkg/net/http/#Response
	// https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	defer func() {
		io.Copy(io.Discard, resp.Body)
	}()

	// Don't try to decode on 204s or Content-Length is 0
	if resp.StatusCode == http.StatusNoContent || resp.ContentLength == 0 {
		return resp, nil
	}

	// Decode from json
	if successV != nil || failureV != nil {
		err = decodeResponse(resp, s.responseDecoder, successV, failureV)
	}
	return resp, err
}

// decodeResponse decodes response Body into the value pointed to by successV
// if the response is a success (2XX) or into the value pointed to by failureV
// otherwise. If the successV or failureV argument to decode into is nil,
// decoding is skipped.
// Caller is responsible for closing the resp.Body.
func decodeResponse(resp *http.Response, decoder ResponseDecoder, successV, failureV interface{}) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if successV != nil {
			return decoder.Decode(resp, successV)
		}
	} else {
		if failureV != nil {
			return decoder.Decode(resp, failureV)
		}
	}
	return nil
}
