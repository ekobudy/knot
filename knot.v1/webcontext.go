package knot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eaciit/toolkit"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

var (
	DefaultOutputType OutputType
)

type WebContext struct {
	Config  *ResponseConfig
	Server  *Server
	Request *http.Request
	Writer  http.ResponseWriter

	queryKeys []string
	cookies   map[string]*http.Cookie
}

func (r *WebContext) QueryKeys() []string {
	if len(r.queryKeys) == 0 {
		if r.Request == nil {
			return r.queryKeys
		}

		values := r.Request.URL.Query()
		for k, _ := range values {
			r.queryKeys = append(r.queryKeys, k)
		}
	}
	return r.queryKeys
}

func (r *WebContext) Query(id string) string {
	if r.Request == nil {
		return ""
	}

	return r.Request.URL.Query().Get(id)
}

func (r *WebContext) Form(id string) string {
	if r.Request == nil {
		return ""
	}
	return r.Request.FormValue(id)
}

func (r *WebContext) GetPayload(result interface{}) error {
	if r.Request == nil {
		return errors.New("HttpRequest object is not properly setup")
	}

	bs, e := ioutil.ReadAll(r.Request.Body)
	if e != nil {
		return fmt.Errorf("Unable to read body: " + e.Error())
	}
	defer r.Request.Body.Close()

	br := bytes.NewReader(bs)
	decoder := json.NewDecoder(br)
	edecode := decoder.Decode(result)
	if edecode != nil {
		return fmt.Errorf("Payload Decode Error: " + edecode.Error() + " .Bytes Data: " + string(bs))
	} else {
		return nil
	}
}

func (r *WebContext) GetForms(result interface{}) error {
	if r.Request == nil {
		return errors.New("HttpRequest object is not properly setup")
	}

	m := toolkit.M{}
	e := r.Request.ParseForm()
	if e != nil {
		return e
	}
	for k, v := range r.Request.Form {
		//fmt.Println("Receiving form %s : %v \n", k, v)
		if f, floatOk := toolkit.StringToFloat(v[0]); floatOk {
			m.Set(k, f)
		} else {
			m.Set(k, v[0])
		}
	}
	e = toolkit.Unjson(m.ToBytes("json", nil), result)
	return e
}

func (r *WebContext) GetPayloadMultipart(result interface{}) (map[string][]*multipart.FileHeader,
	map[string][]string, error) {
	var e error
	if r.Request == nil {
		return nil, nil, errors.New("HttpRequest object is not properly setup")
	}
	e = r.Request.ParseMultipartForm(1024 * 1024 * 1024 * 1024)
	if e != nil {
		return nil, nil, fmt.Errorf("Unable to parse: %s", e.Error())
	}
	m := r.Request.MultipartForm
	return m.File, m.Value, nil
}
