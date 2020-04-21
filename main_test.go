package swagger

import (
	"net/http"
	"testing"

	_ "github.com/arsmn/fiber-swagger/example/docs"
	"github.com/gofiber/fiber"
)

func Test_Swagger(t *testing.T) {
	app := *fiber.New()

	app.Use(New())

	tests := []struct {
		name        string
		url         string
		statusCode  int
		contentType string
	}{
		{
			name:        "Should be returns status 200 with 'text/html; charset=utf-8' content-type",
			url:         "/swagger/index.html",
			statusCode:  200,
			contentType: "text/html; charset=utf-8",
		},
		{
			name:        "Should be returns status 200 with 'application/json; charset=utf-8' content-type",
			url:         "/swagger/doc.json",
			statusCode:  200,
			contentType: "application/json; charset=utf-8",
		},
		{
			name:        "Should be returns status 200 with 'image/png' content-type",
			url:         "/swagger/favicon-16x16.png",
			statusCode:  200,
			contentType: "image/png",
		},
		{
			name:       "Should return status 404",
			url:        "/swagger/notfound",
			statusCode: 404,
		},
		{
			name:       "Should return status 301",
			url:        "/swagger",
			statusCode: 301,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf(`%s: %s`, t.Name(), err)
			}

			if resp.StatusCode != tt.statusCode {
				t.Fatalf(`%s: StatusCode: got %v - expected %v`, t.Name(), resp.StatusCode, tt.statusCode)
			}

			if tt.contentType != "" {
				ct := resp.Header.Get("Content-Type")
				if ct != tt.contentType {
					t.Fatalf(`%s: Content-Type: got %s - expected %s`, t.Name(), ct, tt.contentType)
				}
			}
		})
	}

}
