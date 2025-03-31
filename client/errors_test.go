package client

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestServerSideError(t *testing.T) {
	t.Run("response", func(t *testing.T) {
		expected := &ServerSideError{
			Code:       "my-code",
			Message:    "my-message",
			RequestID:  "my-request-id",
			HostID:     "my-host-id",
			StatusCode: fasthttp.StatusBadRequest,
		}

		resp := new(fasthttp.Response)
		resp.SetStatusCode(expected.StatusCode)
		resp.SetBody([]byte(`<Error>
			<Code>` + expected.Code + `</Code>
			<Message>` + expected.Message + `</Message>
			<RequestId>` + expected.RequestID + `</RequestId>
			<HostId>` + expected.HostID + `</HostId>
		</Error>`))

		actual := NewServerSideError(resp)
		actual.Response = nil

		require.Equal(t, expected, actual)
	})

	t.Run("message", func(t *testing.T) {
		sse := &ServerSideError{
			Code:      "my-code",
			Message:   "my-message",
			RequestID: "my-request-id",
			HostID:    "my-host-id",
		}

		require.Equal(
			t,
			"server-side error occurred (ErrorCode:my-code) (RequestID:my-request-id) (HostID:my-host-id): my-message",
			sse.Error(),
		)
	})
}
