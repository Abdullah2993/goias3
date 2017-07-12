package ias3

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const baseURL = "http://s3.us.archive.org/"

//Req holds the data for the request
type Req struct {
	header http.Header
}

//NewReq returs a new Req
func NewReq() *Req {
	return &Req{header: make(http.Header)}
}

//AutoCreateBucket creates a bucket automatically when uploading a file
func (r *Req) AutoCreateBucket() *Req {
	r.header.Add("x-archive-auto-make-bucket", "1")
	return r
}

//WithCreds sets the credentials for the Req
func (r *Req) WithCreds(accessKey, secretKey string) *Req {
	r.header.Add("Authorization", fmt.Sprintf("LOW %s:%s", accessKey, secretKey))
	return r
}

//WithMeta sets the metadata for the bucket or file
func (r *Req) WithMeta(name, value string) *Req {
	name = strings.Replace(name, "_", "--", -1)
	r.header.Add("x-archive-meta-"+name, value)
	return r
}

//WithMetaMulti sets the metadata for the bucket or file with multiple tags of same name
func (r *Req) WithMetaMulti(name string, value ...string) *Req {
	name = strings.Replace(name, "_", "--", -1)
	for i, val := range value {
		i++
		r.header.Add(fmt.Sprintf("x-archive-meta%d-%s", i, name), val)
	}
	return r
}

//SkipDerive skips the derive process
func (r *Req) SkipDerive() *Req {
	r.header.Add("x-archive-queue-derive", "0")
	return r
}

//KeepOldVersion keeps the old version of the file if it exists by default the
//old file is overridden
func (r *Req) KeepOldVersion() *Req {
	r.header.Add("x-archive-keep-old-version", "1")
	return r
}

//Interactive sets the priority to high so the file is processed imediatley
func (r *Req) Interactive() *Req {
	r.header.Add("x-archive-interactive-priority", "1")
	return r
}

//CreateBucket creates a new bucket
func (r *Req) CreateBucket(key string) error {
	if validateBucketKey(key) {
		return fmt.Errorf("invalid bucket key")
	}
	if strings.Contains(key, "/") {
		return fmt.Errorf("key contains illegal charaacters")
	}
	req, err := http.NewRequest("PUT", baseURL+key, nil)
	if err != nil {
		return err
	}

	return checkResp(http.DefaultClient.Do(req))

}

//UpdateBucket updates an existing bucket
func (r *Req) UpdateBucket(key string, ignorePrevious bool) error {
	if validateBucketKey(key) {
		return fmt.Errorf("invalid bucket key")
	}
	req, err := http.NewRequest(http.MethodPut, baseURL+key, nil)
	if err != nil {
		return err
	}

	if ignorePrevious {
		req.Header.Add("x-archive-ignore-preexisting-bucket", "1")
	}
	r.copy(req)
	return checkResp(http.DefaultClient.Do(req))
}

//Upload uploads a file to a bucket, for large files use UploadFile
func (r *Req) Upload(key string, data []byte) error {
	if validateFileKey(key) {
		return fmt.Errorf("invalid file key")
	}
	req, err := http.NewRequest(http.MethodPut, baseURL+key, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	r.copy(req)

	return checkResp(http.DefaultClient.Do(req))
}

//DeleteFile Deletes a file from a bucket. Internet Archive doesn't
//support deletion of buckets
func (r *Req) DeleteFile(key string) error {
	if validateFileKey(key) {
		return fmt.Errorf("invalid File key")
	}
	req, err := http.NewRequest(http.MethodDelete, baseURL+key, nil)
	if err != nil {
		return err
	}
	r.copy(req)

	return checkResp(http.DefaultClient.Do(req))
}

//UploadFile uploads a file, and closes the input file
func (r *Req) UploadFile(key string, file *os.File) error {
	if validateFileKey(key) {
		return fmt.Errorf("invalid file key")

	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file stats: %v", err)
	}
	req, err := http.NewRequest(http.MethodPut, baseURL+key, file)
	r.copy(req)
	req.ContentLength = fi.Size()
	return checkResp(http.DefaultClient.Do(req))
}

//copy copies headers from Req to *http.Request
func (r *Req) copy(req *http.Request) {
	for h, vs := range r.header {
		req.Header[h] = vs
	}
}

func checkResp(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// var msg string
		// b, err := ioutil.ReadAll(resp.Body)
		// if err == nil {
		// 	msg = string(b)
		// }
		return fmt.Errorf("failed with %d: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func validateBucketKey(key string) bool {
	m, _ := regexp.Match("^[A-Za-z0-9\\-_]+$", []byte(key))
	return m
}

func validateFileKey(key string) bool {
	m := false
	if strings.HasPrefix(key, "/") {
		return m
	}
	m, _ = regexp.Match("^[A-Za-z0-9\\-_\\/]+$", []byte(key))
	return m
}
