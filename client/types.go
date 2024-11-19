package client

import "encoding/xml"

type CreateBucketConfiguration struct {
	XMLName            xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CreateBucketConfiguration"`
	LocationConstraint *LocationConstraint
}

type LocationConstraint string
