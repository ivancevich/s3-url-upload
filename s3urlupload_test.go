package s3urlupload

import (
	"strings"
	"testing"
)

func Test(t *testing.T) {
	media := []string{
		"http://www.eventelectronics.com/images/2020-gallery-001.jpg",
		"http://www.eventelectronics.com/images/2020-gallery-002.jpg",
		"http://www.eventelectronics.com/images/2020-gallery-003.jpg",
		"http://www.eventelectronics.com/images/2020-gallery-004.jpg",
	}

	config := Config{
		AwsAccessKey:  "{AWS_ACCESS_KEY}",
		AwsSecretKey:  "{AWS_SECRET_KEY}",
		AwsS3Endpoint: "{AWS_S3_ENDPOINT}", // e.g. "s3-us-west-2.amazonaws.com"
		AwsS3Bucket:   "{AWS_S3_BUCKET}",
		Workers:       4,
		GetFilePath: func(url string) string {
			parts := strings.Split(url, "/")
			return "eventelectronics/" + parts[len(parts)-1]
		},
	}

	s3UrlUpload := Init(config)
	results := s3UrlUpload.Run(media...)

	for r := range results {
		if r.Error != nil {
			t.Fail()
		}
	}
}
