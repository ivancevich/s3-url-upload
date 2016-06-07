package s3urlupload

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/rlmcpherson/s3gof3r"
)

type Config struct {
	AwsAccessKey  string
	AwsSecretKey  string
	AwsS3Endpoint string
	AwsS3Bucket   string
	Workers       uint
	GetFilePath   func(string) string
}

func Init(c Config) *S3UrlUpload {
	if c.Workers == 0 {
		c.Workers = 1
	}

	if c.GetFilePath == nil {
		c.GetFilePath = func(url string) string {
			parts := strings.Split(url, "/")
			return parts[len(parts)-1]
		}
	}

	s3 := s3gof3r.New(c.AwsS3Endpoint, s3gof3r.Keys{
		AccessKey: c.AwsAccessKey,
		SecretKey: c.AwsSecretKey,
	})

	b := s3.Bucket(c.AwsS3Bucket)

	return &S3UrlUpload{
		config: &c,
		bucket: b,
	}
}

type download struct {
	Body  io.ReadCloser
	URL   string
	Name  string
	Error error
}

type S3UrlUpload struct {
	config *Config
	bucket *s3gof3r.Bucket
}

func (s3uu *S3UrlUpload) Run(files ...string) []error {
	count := len(files)
	jobs := make(chan string, count)
	results := make(chan error, count)

	for w := 1; w <= int(s3uu.config.Workers); w++ {
		go s3uu.worker(jobs, results)
	}

	for _, f := range files {
		jobs <- f
	}

	close(jobs)

	errors := []error{}

	for i := 0; i < count; i++ {
		err := <-results
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (s3uu *S3UrlUpload) download(url string) <-chan download {
	out := make(chan download)
	go func() {
		defer close(out)

		d := download{
			URL:  url,
			Name: s3uu.config.GetFilePath(url),
		}

		resp, err := http.Get(url)
		if err != nil {
			d.Error = err
			out <- d
			return
		}

		if resp.Status != "200 OK" {
			d.Error = errors.New("Status was not OK")
			out <- d
			return
		}

		d.Body = resp.Body
		out <- d
	}()
	return out
}

func (s3uu *S3UrlUpload) upload(in <-chan download) <-chan error {
	var doUpload = func(d download) error {
		defer d.Body.Close()

		w, err := s3uu.bucket.PutWriter(d.Name, nil, nil)
		if err != nil {
			return err
		}

		defer w.Close()

		_, err = io.Copy(w, d.Body)
		if err != nil {
			return err
		}

		return nil
	}

	out := make(chan error)
	go func() {
		defer close(out)
		for d := range in {
			if d.Error != nil {
				out <- d.Error
				continue
			}

			out <- doUpload(d)
		}
	}()
	return out
}

func (s3uu *S3UrlUpload) worker(jobs <-chan string, results chan<- error) {
	for j := range jobs {
		err := <-s3uu.upload(s3uu.download(j))
		results <- err
	}
}
