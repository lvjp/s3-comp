package v4

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/lvjp/s3-comp/client/pkg/signer"
	"github.com/lvjp/s3-comp/client/pkg/signer/utils"

	"github.com/valyala/fasthttp"
)

const signatureValueLen = 64
const chunkDataSize = 64 * 1024

var crlf = []byte{'\r', '\n'}
var trailerKeyValueSeparator = []byte{':'}

var emptyHash = utils.Hex(utils.SHA256Hash(nil))

type payloadTransformer interface {
	Prepare() error
	Transform(seed, date, scope string, signingKey []byte)
}

var _ payloadTransformer = (*plainTransformer)(nil)

type plainTransformer struct {
	req      *fasthttp.Request
	signBody bool
}

func (t *plainTransformer) Prepare() error {
	v := "UNSIGNED-PAYLOAD"
	if t.signBody {
		v = utils.Hex(utils.SHA256Hash(t.req.Body()))
	}
	t.req.Header.Set(HeaderXAmzContentSHA256, v)

	return nil
}

func (t *plainTransformer) Transform(_, _, _ string, _ []byte) {}

type streamTransformer struct {
	req          *fasthttp.Request
	buf          bytes.Buffer
	signBody     bool
	trailerName  []byte
	trailerValue []byte
}

func (t *streamTransformer) Prepare() error {
	contentEncoding := "aws-chunked"
	if actual := t.req.Header.ContentEncoding(); len(actual) > 0 {
		contentEncoding += "," + string(actual)
	}

	newLen := t.transformedBodyLen()
	t.buf.Grow(newLen)
	t.req.Header.DelBytes(t.trailerName)
	t.req.Header.Set(HeaderXAmzContentSHA256, t.contentSHA256())
	t.req.Header.SetContentEncoding(contentEncoding)
	t.req.Header.Set(HeaderXAmzDecodedContentLength, strconv.Itoa(len(t.req.Body())))
	t.req.Header.SetContentLength(newLen)

	return nil
}

func (t *streamTransformer) contentSHA256() string {
	if !t.signBody {
		return "STREAMING-UNSIGNED-PAYLOAD-TRAILER"
	}

	if len(t.trailerName) == 0 {
		return "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"
	}

	return "STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER"
}

func (t *streamTransformer) transformedBodyLen() int {
	const signatureSize = len(";chunk-signature=") + signatureValueLen

	decodedContentLength := len(t.req.Body())

	bodyLen := decodedContentLength
	if nbChunk := decodedContentLength / chunkDataSize; nbChunk > 0 {
		chunkSize := len(strconv.FormatInt(chunkDataSize, 16)) + len(crlf)
		if t.signBody {
			chunkSize += signatureSize
		}
		chunkSize += len(crlf)

		bodyLen += nbChunk * chunkSize
	}

	if remaining := decodedContentLength % chunkDataSize; remaining > 0 {
		bodyLen += len(strconv.FormatInt(int64(remaining), 16)) + len(crlf)
		if t.signBody {
			bodyLen += signatureSize
		}
		bodyLen += len(crlf)
	}

	bodyLen += len("0") + len(crlf)
	if t.signBody {
		bodyLen += signatureSize
	}

	if t.trailerName != nil {
		bodyLen += len(t.trailerName) + len(trailerKeyValueSeparator) + len(t.trailerValue) + len(crlf)

		if t.signBody {
			bodyLen += len("x-amz-trailer-signature:") + signatureValueLen + len(crlf)
		}

		bodyLen += len(crlf)
	}

	return bodyLen
}

func (t *streamTransformer) Transform(seed, date, scope string, signingKey []byte) {
	writeChunk := func(chunk []byte, previousSignature string) string {
		stringToSign := fmt.Sprintf(
			"AWS4-HMAC-SHA256-PAYLOAD\n%s\n%s\n%s\n%s\n%s",
			date,
			scope,
			previousSignature,
			emptyHash,
			utils.Hex(utils.SHA256Hash(chunk)),
		)

		currentSignature := utils.Hex(utils.HMAC_SHA256(signingKey, []byte(stringToSign)))

		t.buf.WriteString(strconv.FormatInt(int64(len(chunk)), 16))
		if t.signBody {
			t.buf.WriteString(";chunk-signature=")
			t.buf.WriteString(currentSignature)
		}
		t.buf.WriteString("\r\n")
		if len(chunk) > 0 {
			t.buf.Write(chunk)
			t.buf.WriteString("\r\n")
		}

		return currentSignature
	}

	for chunk := range slices.Chunk(t.req.Body(), chunkDataSize) {
		seed = writeChunk(chunk, seed)
	}

	seed = writeChunk(nil, seed)

	if len(t.trailerName) > 0 {
		stringToSign := fmt.Sprintf(
			"AWS4-HMAC-SHA256-TRAILER\n%s\n%s\n%s\n%s",
			date,
			scope,
			seed,
			utils.Hex(utils.SHA256Hash(slices.Concat(t.trailerName, trailerKeyValueSeparator, t.trailerValue, []byte{'\n'}))),
		)

		currentSignature := utils.Hex(utils.HMAC_SHA256(signingKey, []byte(stringToSign)))
		t.buf.Write(t.trailerName)
		t.buf.Write(trailerKeyValueSeparator)
		t.buf.Write(t.trailerValue)
		t.buf.WriteString("\r\n")
		if t.signBody {
			t.buf.WriteString("x-amz-trailer-signature:")
			t.buf.WriteString(currentSignature)
			t.buf.WriteString("\r\n")
		}
		t.buf.WriteString("\r\n")
	}

	t.req.Header.SetContentLength(t.buf.Len())
	t.req.SetBody(t.buf.Bytes())
}

type HeaderSigner struct {
	signBody       bool
	forceStreaming bool
}

func NewHeaderSigner(signBody, forceStreaming bool) *HeaderSigner {
	return &HeaderSigner{
		signBody:       signBody,
		forceStreaming: forceStreaming,
	}
}

func (s *HeaderSigner) getPayloadTransformer(req *fasthttp.Request) (payloadTransformer, error) {
	trailerName := req.Header.Peek(HeaderXAmzTrailer)
	var trailerValue []byte

	haveTrailer := len(trailerName) > 0

	if haveTrailer {
		trailerValue = req.Header.PeekBytes(trailerName)
		if len(trailerValue) == 0 {
			return nil, fmt.Errorf("no value set for trailer: %q", trailerName)
		}
	}

	if !s.signBody && s.forceStreaming && !haveTrailer {
		return nil, errors.New("cannot stream an unsigned payload without trailer")
	}

	if !s.forceStreaming && !haveTrailer {
		return &plainTransformer{
				req:      req,
				signBody: s.signBody,
			},
			nil
	}

	return &streamTransformer{
			req:          req,
			signBody:     s.signBody,
			trailerName:  trailerName,
			trailerValue: trailerValue,
		},
		nil
}

func (s *HeaderSigner) Sign(req *fasthttp.Request, credentials *signer.Credentials, region string, now time.Time) (canonicalRequest string, stringToSign string, err error) {
	// AWS S3 specify the use of UTC time
	now = now.UTC()

	ctx := &headerSigningCtx{
		Signer: *s,

		Req:         req,
		Credentials: *credentials,
		Region:      region,
		Now:         now.UTC(),
		haveTrailer: len(req.Header.Peek(HeaderXAmzTrailer)) > 0,

		scope: now.Format(FormatYYYYMMDD) + "/" + region + "/s3/aws4_request",
	}

	payloadTransformer, err := s.getPayloadTransformer(req)
	if err != nil {
		return "", "", err
	}

	// Sanitize request
	req.Header.Del(HeaderAuthorization)
	req.Header.Del(HeaderXAmzContentSHA256)
	req.Header.Set(HeaderXAmzDate, now.Format(FormatXAmzDate))

	if err := payloadTransformer.Prepare(); err != nil {
		return "", "", err
	}

	ctx.computeHeaders()

	ctx.computeSigningKey()

	ctx.computeCanonicalRequest()
	ctx.computeStringToSign()

	signature := ctx.getSignature()

	req.Header.Set(
		HeaderAuthorization,
		fmt.Sprintf(
			"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
			credentials.AccessKeyID,
			ctx.scope,
			ctx.signedHeaders,
			signature,
		),
	)

	payloadTransformer.Transform(
		signature,
		ctx.Now.Format(FormatXAmzDate),
		ctx.scope,
		ctx.signingKey,
	)

	return ctx.canonicalRequest, ctx.stringToSign, nil
}

type headerSigningCtx struct {
	Signer HeaderSigner

	Req         *fasthttp.Request
	Credentials signer.Credentials
	Region      string
	Now         time.Time

	haveTrailer bool

	scope            string
	canonicalHeaders string
	signedHeaders    string
	signingKey       []byte

	canonicalRequest string
	stringToSign     string
}

func (ctx *headerSigningCtx) computeSigningKey() {
	dateKey := utils.HMAC_SHA256([]byte("AWS4"+ctx.Credentials.SecretAccessKey), []byte(ctx.Now.Format(FormatYYYYMMDD)))
	dateRegionKey := utils.HMAC_SHA256(dateKey, []byte(ctx.Region))
	dateRegionServiceKey := utils.HMAC_SHA256(dateRegionKey, []byte("s3"))
	ctx.signingKey = utils.HMAC_SHA256(dateRegionServiceKey, []byte("aws4_request"))
}

func (ctx *headerSigningCtx) computeHeaders() {
	normalized := make(map[string]string, ctx.Req.Header.Len())

	ctx.Req.Header.VisitAll(func(key, value []byte) {
		normalizedKey := utils.LowerCase(string(key))
		normalizedValue := utils.Trim(string(value))

		normalized[normalizedKey] = normalizedValue
	})

	// Ensure host header presence
	if _, exists := normalized["host"]; !exists {
		normalized["host"] = string(ctx.Req.Host())
	}

	sortedHeaders := slices.Sorted(maps.Keys(normalized))
	for _, key := range sortedHeaders {
		ctx.canonicalHeaders += key + ":" + normalized[key] + "\n"
	}

	ctx.signedHeaders = strings.Join(sortedHeaders, ";")
}

func (ctx *headerSigningCtx) getSignature() string {
	return utils.Hex(utils.HMAC_SHA256(ctx.signingKey, []byte(ctx.stringToSign)))
}

func (ctx *headerSigningCtx) computeStringToSign() {
	ctx.stringToSign = "AWS4-HMAC-SHA256\n" +
		ctx.Now.Format(FormatXAmzDate) + "\n" +
		ctx.scope + "\n" +
		utils.Hex(utils.SHA256Hash([]byte(ctx.canonicalRequest)))
}

func (ctx *headerSigningCtx) computeCanonicalRequest() {
	var ret strings.Builder

	ret.Write(ctx.Req.Header.Method())
	ret.WriteByte('\n')

	path := string(ctx.Req.URI().PathOriginal())
	if path == "" {
		path = "/"
	}
	ret.WriteString(utils.URIEncode(path, true))
	ret.WriteByte('\n')

	ret.WriteString(ctx.getCanonicalQueryString())
	ret.WriteByte('\n')

	ret.WriteString(ctx.canonicalHeaders)
	ret.WriteByte('\n')

	ret.WriteString(ctx.signedHeaders)
	ret.WriteByte('\n')

	ret.Write(ctx.Req.Header.Peek(HeaderXAmzContentSHA256))

	ctx.canonicalRequest = ret.String()
}

func (ctx *headerSigningCtx) getCanonicalQueryString() string {
	args := ctx.Req.URI().QueryArgs()
	encoded := make(map[string]string, args.Len())

	args.VisitAll(func(key, value []byte) {
		encoded[utils.URIEncode(string(key), false)] = utils.URIEncode(string(value), false)
	})

	var ret strings.Builder

	for i, key := range slices.Sorted(maps.Keys(encoded)) {
		if i > 0 {
			ret.WriteRune('&')
		}

		ret.WriteString(key)
		ret.WriteRune('=')
		ret.WriteString(encoded[key])
	}

	return ret.String()
}
