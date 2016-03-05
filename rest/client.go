package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Client struct {
	Client *http.Client
	Log    bool // Log request and response

	// Optional
	UserInfo *url.Userinfo

	// Optional defaults - can be overridden in a Request
	Header *http.Header
	Params *url.Values
}

// Send constructs and sends an HTTP request.
func (c *Client) Send(r *Request) (response *Response, err error) {
	r.Method = strings.ToUpper(r.Method)

	// Create a URL object from the raw url string.  This will allow us to compose
	// query parameters programmatically and be guaranteed of a well-formed URL.
	u, err := url.Parse(r.Url)
	if err != nil {
		log.Errorf("URL %s, %s", r.Url, err)
		return
	}

	// Default query parameters
	p := url.Values{}
	if c.Params != nil {
		for k, v := range *c.Params {
			p[k] = v
		}
	}

	// Parameters that were present in URL
	if u.Query() != nil {
		for k, v := range u.Query() {
			p[k] = v
		}
	}

	// User-supplied params override default
	if r.Params != nil {
		for k, v := range *r.Params {
			p[k] = v
		}
	}
	//
	// Encode parameters
	//
	u.RawQuery = p.Encode()
	//
	// Attach params to response
	//
	r.Params = &p
	//
	// Create a Request object; if populated, Data field is JSON encoded as
	// request body
	//
	header := http.Header{}
	if c.Header != nil {
		for k, _ := range *c.Header {
			v := c.Header.Get(k)
			header.Set(k, v)
		}
	}
	var req *http.Request
	var buf *bytes.Buffer
	if r.Payload != nil {
		if r.RawPayload {
			var ok bool
			// buf can be nil interface at this point
			// so we'll do extra nil check
			buf, ok = r.Payload.(*bytes.Buffer)
			if !ok {
				err = errors.New("Payload must be of type *bytes.Buffer if RawPayload is set to true")
				return
			}
		} else {
			var b []byte
			b, err = json.Marshal(&r.Payload)
			if err != nil {
				log.Errorln(err)
				return
			}
			buf = bytes.NewBuffer(b)
		}
		if buf != nil {
			req, err = http.NewRequest(r.Method, u.String(), buf)
		} else {
			req, err = http.NewRequest(r.Method, u.String(), nil)
		}
		if err != nil {
			log.Errorln(err)
			return
		}
		// Overwrite the content type to json since we're pushing the payload as json
		header.Set("Content-Type", "application/json")
	} else { // no data to encode
		req, err = http.NewRequest(r.Method, u.String(), nil)
		if err != nil {
			log.Errorln(err)
			return
		}

	}
	//
	// Merge Session and Request options
	//
	var userinfo *url.Userinfo
	if u.User != nil {
		userinfo = u.User
	}
	if c.UserInfo != nil {
		userinfo = c.UserInfo
	}

	header.Add("Accept", "application/json")
	req.Header = header

	//
	// Set HTTP Basic authentication if userinfo is supplied
	//
	if userinfo != nil {
		pwd, _ := userinfo.Password()
		req.SetBasicAuth(userinfo.Username(), pwd)
		if u.Scheme != "https" {
			log.Debugln("WARNING: Using HTTP Basic Auth in cleartext is insecure.")
		}
	}
	//
	// Execute the HTTP request
	//

	// Debug log request
	log.Debugln("--------------------------------------------------------------------------------")
	log.Debugln("REQUEST")
	log.Debugln("--------------------------------------------------------------------------------")
	log.Debugln("Method:", req.Method)
	log.Debugln("URL:", req.URL)
	log.Debugln("Header:", req.Header)
	log.Debugln("Form:", req.Form)
	log.Debugln("Payload:")
	if r.RawPayload && c.Log && buf != nil {
		log.Debugln(base64.StdEncoding.EncodeToString(buf.Bytes()))
	} else {
		log.Debugln(pretty(r.Payload))
	}
	r.timestamp = time.Now()
	var client *http.Client
	if c.Client != nil {
		client = c.Client
	} else {
		client = &http.Client{}
		c.Client = client
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer resp.Body.Close()
	r.status = resp.StatusCode
	r.response = resp

	//
	// Unmarshal
	//
	r.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return
	}
	if string(r.body) != "" {
		if resp.StatusCode < 300 && r.Result != nil {
			err = json.Unmarshal(r.body, r.Result)
		}
		if resp.StatusCode >= 400 && r.Error != nil {
			json.Unmarshal(r.body, r.Error) // Should we ignore unmarshall error?
		}
	}
	if r.CaptureResponseBody {
		r.ResponseBody = bytes.NewBuffer(r.body)
	}
	rsp := Response(*r)
	response = &rsp

	// Debug log response
	log.Debugln("--------------------------------------------------------------------------------")
	log.Debugln("RESPONSE")
	log.Debugln("--------------------------------------------------------------------------------")
	log.Debugln("Status: ", response.status)
	log.Debugln("Header:")
	log.Debugln(pretty(response.HttpResponse().Header))
	log.Debugln("Body:")

	if response.body != nil {
		raw := json.RawMessage{}
		if json.Unmarshal(response.body, &raw) == nil {
			log.Debugln(pretty(&raw))
		} else {
			log.Debugln(pretty(response.RawText()))
		}
	} else {
		log.Debugln("Empty response body")
	}

	return
}

// Get sends a GET request.
func (c *Client) Get(url string) (*Response, error) {
	r := Request{
		Method: "GET",
		Url:    url,
	}
	return c.Send(&r)
}

// Post sends a POST request.
func (c *Client) Post(url string, payload interface{}) (*Response, error) {
	r := Request{
		Method:     "POST",
		Url:        url,
		Payload:    payload,
		RawPayload: true,
	}
	return c.Send(&r)
}

// Put sends a PUT request.
func (c *Client) Put(url string, payload interface{}) (*Response, error) {
	r := Request{
		Method:     "PUT",
		Url:        url,
		Payload:    payload,
		RawPayload: true,
	}
	return c.Send(&r)
}

// Delete sends a DELETE request.
func (c *Client) Delete(url string) (*Response, error) {
	r := Request{
		Method: "DELETE",
		Url:    url,
	}
	return c.Send(&r)
}
