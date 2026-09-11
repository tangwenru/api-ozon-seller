package ozon

import (
	"context"
	"net/http"
	"testing"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

func TestGetRoles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode int
		headers    map[string]string
		response   string
	}{
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			`{
				"expires_at": "2026-02-18T09:54:23.296Z",
				"roles": [
					{ "name": "Admin", "methods": ["/v1/actions"] },
					{ "name": "Posting FBS", "methods": ["/v1/posting"] }
				]
			}`,
		},
		{
			http.StatusUnauthorized,
			map[string]string{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Roles().GetRoles(ctx)
		if err != nil {
			t.Error(err)
			continue
		}

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode == http.StatusOK {
			if resp.ExpiresAt.IsZero() {
				t.Errorf("ExpiresAt cannot be zero")
			}
			if len(resp.Roles) == 0 {
				t.Errorf("Roles cannot be empty")
			}
		}
	}
}
