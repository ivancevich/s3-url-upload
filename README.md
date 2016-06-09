# s3-url-upload
Upload files to S3 from URLs in Golang

## Start using it

1) Download and install it:

```bash
$ go get github.com/ivancevich/s3-url-upload
```

2) Import it in your code:

```go
import "github.com/ivancevich/s3-url-upload"
```

## Example

```go
package main

import (
	"fmt"
	"strings"
	"github.com/ivancevich/s3-url-upload"
)

func main() {
	media := []string{
		"http://www.eventelectronics.com/images/2020-gallery-001.jpg",
		"http://www.eventelectronics.com/images/2020-gallery-002.jpg",
		"http://www.eventelectronics.com/images/2020-gallery-003.jpg",
		"http://www.eventelectronics.com/images/2020-gallery-004.jpg",
	}

	config := s3urlupload.Config{
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

	s3UrlUpload := s3urlupload.Init(config)
	results := s3UrlUpload.Run(media...)

	for r := range results {
		fmt.Printf("%+v\n", r)
	}
}
```
