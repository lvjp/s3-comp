package types

type ListAllMyBucketsResult struct {
	Buckets           []Bucket `xml:">Bucket"`
	Owner             *Owner
	ContinuationToken *string
	Prefix            *string
}

type Bucket struct {
	BucketRegion *string
	CreationDate *string
	Name         *string
}

type Owner struct {
	DisplayName *string
	ID          *string
}

type CreateBucketConfiguration struct {
	LocationConstraint *LocationConstraint
}

type LocationConstraint string
