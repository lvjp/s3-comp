package v4

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/lvjp/s3-comp/client/pkg/signer"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestHeaderSigner_Sign(t *testing.T) {
	testCases := []struct {
		name                    string
		signBody                bool
		forceStreaming          bool
		req                     func(*fasthttp.Request)
		expectedCanonicalString []string
		expectedStringToSign    []string
		expectedHeaders         []string
		expectedBody            []string
	}{
		{
			name:           "UNSIGNED-PAYLOAD (no payload)",
			signBody:       false,
			forceStreaming: false,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("HEAD")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/0bcfb63f-d90b-49c0-9cc4-cfe8a375c022")
			},
			expectedCanonicalString: []string{
				"HEAD",
				"/0bcfb63f-d90b-49c0-9cc4-cfe8a375c022",
				"",
				"host:examplebucket.s3.amazonaws.com",
				"x-amz-content-sha256:UNSIGNED-PAYLOAD",
				"x-amz-date:19840805T135000Z",
				"",
				"host;x-amz-content-sha256;x-amz-date",
				"UNSIGNED-PAYLOAD",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"f681c3074c5a4df7256d337d9fb24c7e609a80c60be4b3e05bb705cc923eb7a8",
			},
			expectedHeaders: []string{
				"HEAD /0bcfb63f-d90b-49c0-9cc4-cfe8a375c022 HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: UNSIGNED-PAYLOAD",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=c46062ede51e98507b13d77ae00314cf9253b9812cd00927fb30480e8db05d22",
			},
		},
		{
			name: "UNSIGNED-PAYLOAD (with payload)",

			signBody:       false,
			forceStreaming: false,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("PUT")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/test.txt")
				req.Header.Set("X-Amz-Checksum-Crc64nvme", "ntuPBsmdl18=")
				req.Header.SetContentLength(14)
				req.SetBody([]byte("Welcome to S3."))
			},
			expectedCanonicalString: []string{
				"PUT",
				"/test.txt",
				"",
				"content-length:14",
				"host:examplebucket.s3.amazonaws.com",
				"x-amz-checksum-crc64nvme:ntuPBsmdl18=",
				"x-amz-content-sha256:UNSIGNED-PAYLOAD",
				"x-amz-date:19840805T135000Z",
				"",
				"content-length;host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date",
				"UNSIGNED-PAYLOAD",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"c2d8ab094d8d2fa65f3ef72b80f83bf6d4d1a495d7b712d96fc21c14d078c6b2",
			},
			expectedHeaders: []string{
				"PUT /test.txt HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"Content-Length: 14",
				"X-Amz-Checksum-Crc64nvme: ntuPBsmdl18=",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: UNSIGNED-PAYLOAD",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=content-length;host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date, Signature=a9ee21c83da45c070bd9815588c2f04983ec46aa7a5826b2bea5b20b44933cef",
			},
			expectedBody: []string{
				"Welcome to S3.",
			},
		},
		{
			name:           "SIGNED-PAYLOAD (no payload)",
			signBody:       true,
			forceStreaming: false,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("GET")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/test.txt")
				req.Header.Set("Range", "bytes=0-9")
			},
			expectedCanonicalString: []string{
				"GET",
				"/test.txt",
				"",
				"host:examplebucket.s3.amazonaws.com",
				"range:bytes=0-9",
				"x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"x-amz-date:19840805T135000Z",
				"",
				"host;range;x-amz-content-sha256;x-amz-date",
				"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"765e5b5c7ecb1514445224b6dc57b50bad96beda84781d026b596b203d88535b",
			},
			expectedHeaders: []string{
				"GET /test.txt HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"Range: bytes=0-9",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=host;range;x-amz-content-sha256;x-amz-date, Signature=e60941d2d7d9cf5c04cfa3670b551c1defa95cd2e0bd3028674cca38109bdf22",
			},
		},
		{
			name:           "SIGNED-PAYLOAD (with payload)",
			signBody:       true,
			forceStreaming: false,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("PUT")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/test.txt?x-id=PutObject")
				req.Header.Set("X-Amz-Checksum-Crc64nvme", "ntuPBsmdl18=")
				req.Header.SetContentLength(14)
				req.SetBody([]byte("Welcome to S3."))
			},
			expectedCanonicalString: []string{
				"PUT",
				"/test.txt",
				"x-id=PutObject",
				"content-length:14",
				"host:examplebucket.s3.amazonaws.com",
				"x-amz-checksum-crc64nvme:ntuPBsmdl18=",
				"x-amz-content-sha256:f3893d4cc3e907c99afd2b35ae83e391b914b78c98097d9b5f7c89d4800fbaa9",
				"x-amz-date:19840805T135000Z",
				"",
				"content-length;host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date",
				"f3893d4cc3e907c99afd2b35ae83e391b914b78c98097d9b5f7c89d4800fbaa9",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"f27e60ea6b7505eb7fb41bd1f491987aa0e90b06b43fc00f14a107664293c754",
			},
			expectedHeaders: []string{
				"PUT /test.txt?x-id=PutObject HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"Content-Length: 14",
				"X-Amz-Checksum-Crc64nvme: ntuPBsmdl18=",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: f3893d4cc3e907c99afd2b35ae83e391b914b78c98097d9b5f7c89d4800fbaa9",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=content-length;host;x-amz-checksum-crc64nvme;x-amz-content-sha256;x-amz-date, Signature=7aa44e95e43973edb9b2af6fff6461e92bdc623921f17d2ac19a1757f3cd06fe",
			},
			expectedBody: []string{"Welcome to S3."},
		},
		{
			name:           "SIGNED STREAMING",
			signBody:       true,
			forceStreaming: true,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("PUT")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/a02f63c7-3841-4e3d-8e51-87d7f80ce655")
				req.Header.SetContentLength(66560)
				req.SetBody(slices.Repeat([]byte{'a'}, 65*1024))
			},
			expectedCanonicalString: []string{
				"PUT",
				"/a02f63c7-3841-4e3d-8e51-87d7f80ce655",
				"",
				"content-encoding:aws-chunked",
				"content-length:66822",
				"host:examplebucket.s3.amazonaws.com",
				"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
				"x-amz-date:19840805T135000Z",
				"x-amz-decoded-content-length:66560",
				"",
				"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length",
				"STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"1cb4bcdc0f41cdabb7500cd1a5b04d63221cdcc2c7366207f520bcabfea32af3",
			},
			expectedHeaders: []string{
				"PUT /a02f63c7-3841-4e3d-8e51-87d7f80ce655 HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"Content-Length: 66822",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD",
				"Content-Encoding: aws-chunked",
				"X-Amz-Decoded-Content-Length: 66560",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length, Signature=ed20e4eaccc6bf87aefc39a735e505f7c7b8f5a9123256d1e7e305ee6d46811d",
			},
			expectedBody: []string{
				"10000;chunk-signature=5c0e89f79f041ccf739707cf8832397dcc74c2c56435ffd6b2e1d71d5da87c88\r\n",
				strings.Repeat("a", 64*1024) + "\r\n",
				"400;chunk-signature=49b3a44d3ebbd53f810584b40fc01bae82959c8a684117be3e77dbf7e6ddd882\r\n",
				strings.Repeat("a", 1024) + "\r\n",
				"0;chunk-signature=78e649f3c60d2af2c4eaef16482d0574b84355e036abf6fbbeb5387306c6e0a7\r\n",
			},
		},
		{
			name:           "UNSIGNED TRAILER",
			signBody:       false,
			forceStreaming: false,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("PUT")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/test.txt?x-id=PutObject")
				req.Header.Set("X-Amz-Sdk-Checksum-Algorithm", "CRC64NVME")
				req.Header.Set("X-Amz-Trailer", "x-amz-checksum-crc64nvme")
				req.Header.Set("x-amz-checksum-crc64nvme", "ntuPBsmdl18=")
				req.Header.SetContentLength(14)
				req.SetBody([]byte("Welcome to S3."))
			},
			expectedCanonicalString: []string{
				"PUT",
				"/test.txt",
				"x-id=PutObject",
				"content-encoding:aws-chunked",
				"content-length:63",
				"host:examplebucket.s3.amazonaws.com",
				"x-amz-content-sha256:STREAMING-UNSIGNED-PAYLOAD-TRAILER",
				"x-amz-date:19840805T135000Z",
				"x-amz-decoded-content-length:14",
				"x-amz-sdk-checksum-algorithm:CRC64NVME",
				"x-amz-trailer:x-amz-checksum-crc64nvme",
				"",
				"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-sdk-checksum-algorithm;x-amz-trailer",
				"STREAMING-UNSIGNED-PAYLOAD-TRAILER",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"aebdea80272e30161f628e2df45400b07b66cf201c1253abdcf9c050e02feae8",
			},
			expectedHeaders: []string{
				"PUT /test.txt?x-id=PutObject HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"Content-Length: 63",
				"X-Amz-Sdk-Checksum-Algorithm: CRC64NVME",
				"X-Amz-Trailer: x-amz-checksum-crc64nvme",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: STREAMING-UNSIGNED-PAYLOAD-TRAILER",
				"Content-Encoding: aws-chunked",
				"X-Amz-Decoded-Content-Length: 14",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-sdk-checksum-algorithm;x-amz-trailer, Signature=c01dcd51b9d4308d3c68703c39544da77abb9f83bceb2f4d03e7b739cbab5077",
			},
			expectedBody: []string{
				"e\r\nWelcome to S3.\r\n",
				"0\r\n",
				"x-amz-checksum-crc64nvme:ntuPBsmdl18=\r\n",
				"\r\n",
			},
		},
		{
			name:     "SIGNED TRAILER",
			signBody: true,
			req: func(req *fasthttp.Request) {
				req.Header.SetMethod("PUT")
				req.SetRequestURI("https://examplebucket.s3.amazonaws.com/3014120a-9e80-4956-8f84-60c79cb8013f")
				req.Header.Set("X-Amz-Trailer", "x-amz-checksum-crc32")
				req.Header.Set("x-amz-checksum-crc32", "sK4Y7A==")
				req.Header.SetContentLength(65 * 1024)
				req.SetBody(slices.Repeat([]byte{'a'}, 65*1024))
			},
			expectedCanonicalString: []string{
				"PUT",
				"/3014120a-9e80-4956-8f84-60c79cb8013f",
				"",
				"content-encoding:aws-chunked",
				"content-length:66945",
				"host:examplebucket.s3.amazonaws.com",
				"x-amz-content-sha256:STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
				"x-amz-date:19840805T135000Z",
				"x-amz-decoded-content-length:66560",
				"x-amz-trailer:x-amz-checksum-crc32",
				"",
				"content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer",
				"STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
			},
			expectedStringToSign: []string{
				"AWS4-HMAC-SHA256",
				"19840805T135000Z",
				"19840805/eu-west-3/s3/aws4_request",
				"dcb493c343c033658c0e4279d353c3bc33b5213e5fead016c22ed39578851866",
			},
			expectedHeaders: []string{
				"PUT /3014120a-9e80-4956-8f84-60c79cb8013f HTTP/1.1",
				"Host: examplebucket.s3.amazonaws.com",
				"Content-Length: 66945",
				"X-Amz-Trailer: x-amz-checksum-crc32",
				"X-Amz-Date: 19840805T135000Z",
				"X-Amz-Content-Sha256: STREAMING-AWS4-HMAC-SHA256-PAYLOAD-TRAILER",
				"Content-Encoding: aws-chunked",
				"X-Amz-Decoded-Content-Length: 66560",
				"Authorization: AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/19840805/eu-west-3/s3/aws4_request, SignedHeaders=content-encoding;content-length;host;x-amz-content-sha256;x-amz-date;x-amz-decoded-content-length;x-amz-trailer, Signature=bb1c17814c802ebd1a9897a7d9e7e9ac52f304cda6d8a8dd318723808f741f46",
			},
			expectedBody: []string{
				"10000;chunk-signature=7751b9162b6db7f017db55fe50de12328d481426c60a98bb441b1429661a2877\r\n",
				strings.Repeat("a", 64*1024) + "\r\n",
				"400;chunk-signature=fdaef60fb7d99803adf3511f953f18f5a3a2fb3feab3bce60e0a85f760f65af8\r\n",
				strings.Repeat("a", 1024) + "\r\n",
				"0;chunk-signature=3bd5da117d96df0273a0a3eec14c6b627925b32afb5ea8efddd6ce0299d1c399\r\n",
				"x-amz-checksum-crc32:sK4Y7A==\r\n",
				"x-amz-trailer-signature:15afa83816d9377556dd621320b14059885a50066029a635e923492a499aa25d\r\n",
				"\r\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedCanonicalString := strings.Join(tc.expectedCanonicalString, "\n")
			expectedStringToSign := strings.Join(tc.expectedStringToSign, "\n")
			expectedReq := strings.Join(tc.expectedHeaders, "\r\n") + "\r\n\r\n" + strings.Join(tc.expectedBody, "")

			var req fasthttp.Request
			req.Header.SetNoDefaultContentType(true)
			tc.req(&req)

			s := NewHeaderSigner(tc.signBody, tc.forceStreaming)
			actualCannonicalString, actualStringToSign, err := s.Sign(
				&req,
				&signer.Credentials{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
				"eu-west-3",
				time.Date(1984, time.August, 5, 13, 50, 0, 0, time.UTC),
			)
			require.NoError(t, err)

			require.Equal(t, expectedCanonicalString, actualCannonicalString, "canonical string")
			require.Equal(t, expectedStringToSign, actualStringToSign, "string to sign")
			require.Equal(t, expectedReq, req.String(), "request")
		})
	}
}
