package client

import (
	"context"

	"github.com/valyala/fasthttp"
)

var _ RequiredBucketKeyInterface = (*GetObjectInput)(nil)

type GetObjectInput struct {
	// Bucket is mandatory
	Bucket string

	// Key is mandatory
	Key string

	PartNumber                 *string
	ResponseCacheControl       *string
	ResponseContentDisposition *string
	ResponseContentEncoding    *string
	ResponseContentLanguage    *string
	ResponseContentType        *string
	ResponseExpires            *string
	VersionId                  *string

	IfMatch              *string
	IfModifiedSince      *string
	IfNoneMatch          *string
	IfUnmodifiedSince    *string
	Range                *string
	SSECustomerAlgorithm *string
	SSECustomerKey       *string
	SSECustomerKeyMD5    *string
	RequestPayer         *string
	ExpectedBucketOwner  *string
	ChecksumMode         *string
}

func (input *GetObjectInput) GetBucket() string {
	return input.Bucket
}

func (input *GetObjectInput) GetKey() string {
	return input.Key
}

func (input *GetObjectInput) MarshalHTTP(req *fasthttp.Request) error {
	req.Header.SetMethod(fasthttp.MethodGet)

	args := req.URI().QueryArgs()
	setQuery(args, QueryPartNumber, input.PartNumber)
	setQuery(args, QueryResponseCacheControl, input.ResponseCacheControl)
	setQuery(args, QueryResponseContentDisposition, input.ResponseContentDisposition)
	setQuery(args, QueryResponseContentEncoding, input.ResponseContentEncoding)
	setQuery(args, QueryResponseContentLanguage, input.ResponseContentLanguage)
	setQuery(args, QueryResponseContentType, input.ResponseContentType)
	setQuery(args, QueryResponseExpires, input.ResponseExpires)
	setQuery(args, QueryVersionID, input.VersionId)

	setHeader(&req.Header, HeaderIfMatch, input.IfMatch)
	setHeader(&req.Header, HeaderIfModifiedSince, input.IfModifiedSince)
	setHeader(&req.Header, HeaderIfNoneMatch, input.IfNoneMatch)
	setHeader(&req.Header, HeaderIfUnmodifiedSince, input.IfUnmodifiedSince)
	setHeader(&req.Header, HeaderRange, input.Range)
	setHeader(&req.Header, HeaderXAmzSSECustomerAlgorithm, input.SSECustomerAlgorithm)
	setHeader(&req.Header, HeaderXAmzSSECustomerKey, input.SSECustomerKey)
	setHeader(&req.Header, HeaderXAmzSSECustomerKeyMD5, input.SSECustomerKeyMD5)
	setHeader(&req.Header, HeaderXAmzRequestPayer, input.RequestPayer)
	setHeader(&req.Header, HeaderXAmzExpectedBucketOwner, input.ExpectedBucketOwner)
	setHeader(&req.Header, HeaderXAmzChecksumMode, input.ChecksumMode)

	return nil
}

type GetObjectOutput struct {
	Body []byte

	AcceptRanges       *string
	CacheControl       *string
	ContentDisposition *string
	ContentEncoding    *string
	ContentLanguage    *string
	ContentLength      *string
	ContentRange       *string
	ContentType        *string
	ETag               *string
	Expires            *string
	LastModified       *string

	ChecksumCRC32             *string
	ChecksumCRC32C            *string
	ChecksumCRC64NVME         *string
	ChecksumSHA1              *string
	ChecksumSHA256            *string
	ChecksumType              *string
	DeleteMarker              *string
	Expiration                *string
	MissingMeta               *string
	PartsCount                *string
	ObjectLockLegalHoldStatus *string
	ObjectLockMode            *string
	ObjectLockRetainUntilDate *string
	ReplicationStatus         *string
	RequestCharged            *string
	Restore                   *string
	SSEKMSKeyId               *string
	BucketKeyEnabled          *string
	SSECustomerAlgorithm      *string
	SSECustomerKeyMD5         *string
	ServerSideEncryption      *string
	StorageClass              *string
	TaggingCount              *string
	VersionId                 *string
	WebsiteRedirectLocation   *string
}

func (output *GetObjectOutput) UnmarshalHTTP(resp *fasthttp.Response) error {
	if resp.StatusCode() != fasthttp.StatusOK {
		return NewServerSideError(resp)
	}

	output.AcceptRanges = extractHeader(&resp.Header, HeaderAcceptRanges)
	output.CacheControl = extractHeader(&resp.Header, HeaderCacheControl)
	output.ContentDisposition = extractHeader(&resp.Header, HeaderContentDisposition)
	output.ContentEncoding = extractHeader(&resp.Header, HeaderContentEncoding)
	output.ContentLanguage = extractHeader(&resp.Header, HeaderContentLanguage)
	output.ContentLength = extractHeader(&resp.Header, HeaderContentLength)
	output.ContentRange = extractHeader(&resp.Header, HeaderContentRange)
	output.ContentType = extractHeader(&resp.Header, HeaderContentType)
	output.ETag = extractHeader(&resp.Header, HeaderETag)
	output.Expires = extractHeader(&resp.Header, HeaderExpires)
	output.LastModified = extractHeader(&resp.Header, HeaderLastModified)

	output.ChecksumCRC32 = extractHeader(&resp.Header, HeaderXAmzChecksumCRC32)
	output.ChecksumCRC32C = extractHeader(&resp.Header, HeaderXAmzChecksumCRC32C)
	output.ChecksumCRC64NVME = extractHeader(&resp.Header, HeaderXAmzChecksumCRC64NVME)
	output.ChecksumSHA1 = extractHeader(&resp.Header, HeaderXAmzChecksumSHA1)
	output.ChecksumSHA256 = extractHeader(&resp.Header, HeaderXAmzChecksumSHA256)
	output.ChecksumType = extractHeader(&resp.Header, HeaderXAmzChecksumType)
	output.DeleteMarker = extractHeader(&resp.Header, HeaderXAmzDeleteMarker)
	output.Expiration = extractHeader(&resp.Header, HeaderXAmzExpiration)
	output.MissingMeta = extractHeader(&resp.Header, HeaderXAmzMissingMeta)
	output.PartsCount = extractHeader(&resp.Header, HeaderXAmzPartsCount)
	output.ObjectLockLegalHoldStatus = extractHeader(&resp.Header, HeaderXAmzObjectLockLegalHoldStatus)
	output.ObjectLockMode = extractHeader(&resp.Header, HeaderXAmzObjectLockMode)
	output.ObjectLockRetainUntilDate = extractHeader(&resp.Header, HeaderXAmzObjectLockRetainUntilDate)
	output.ReplicationStatus = extractHeader(&resp.Header, HeaderXAmzReplicationStatus)
	output.RequestCharged = extractHeader(&resp.Header, HeaderXAmzRequestCharged)
	output.Restore = extractHeader(&resp.Header, HeaderXAmzRestore)
	output.SSEKMSKeyId = extractHeader(&resp.Header, HeaderXAmzSSEKMSKeyId)
	output.BucketKeyEnabled = extractHeader(&resp.Header, HeaderXAmzBucketKeyEnabled)
	output.SSECustomerAlgorithm = extractHeader(&resp.Header, HeaderXAmzSSECustomerAlgorithm)
	output.SSECustomerKeyMD5 = extractHeader(&resp.Header, HeaderXAmzSSECustomerKeyMD5)
	output.ServerSideEncryption = extractHeader(&resp.Header, HeaderXAmzServerSideEncryption)
	output.StorageClass = extractHeader(&resp.Header, HeaderXAmzStorageClass)
	output.TaggingCount = extractHeader(&resp.Header, HeaderXAmzTaggingCount)
	output.VersionId = extractHeader(&resp.Header, HeaderXAmzVersionId)
	output.WebsiteRedirectLocation = extractHeader(&resp.Header, HeaderXAmzWebsiteRedirectLocation)

	return nil
}

func (c *Client) GetObject(ctx context.Context, input *GetObjectInput, optFns ...func(*Options)) (*GetObjectOutput, *Metadata, error) {
	return PerformCall[*GetObjectInput, *GetObjectOutput](ctx, c, input, optFns...)
}
