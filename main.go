package swagger

import (
	"html/template"
	"path"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/utils"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/swag"
)

const (
	defaultDocURL = "doc.json"
	defaultIndex  = "index.html"
)

// Handler default
var Handler = New()

// Config stores SwaggerUI configuration variables
type Config struct {
	// Enable deep linking for tags and operations, default is true
	DeepLinking bool

	// Controls the default expansion setting for the operations and tags.
	// 'list' (default, expands only the tags),
	// 'full' (expands the tags and operations),
	// 'none' (expands nothing)
	DocExpansion string

	// Configuration information for OAuth2, optional if using OAuth2
	OAuth *OAuthConfig

	// Custom OAuth redirect URL
	OAuth2RedirectUrl string

	// URL pointing to API definition
	URL string
}

type OAuthConfig struct {
	// application name, displayed in authorization popup
	AppName string

	// ID of the client sent to the OAuth2 Provider, default is clientId
	ClientId string
}

// New returns custom handler
func New(config ...Config) fiber.Handler {
	cfg := Config{
		DeepLinking:  true,
		DocExpansion: "list",
	}

	if len(config) > 0 {
		cfg = config[0]
	}

	index, err := template.New("swagger_index.html").Parse(indexTmpl)
	if err != nil {
		panic("swagger: could not parse index template")
	}

	var (
		prefix string
		once   sync.Once
		fs     fiber.Handler = filesystem.New(filesystem.Config{Root: swaggerFiles.HTTP})
	)

	return func(c *fiber.Ctx) error {
		// Set prefix
		once.Do(func() {
			prefix = strings.ReplaceAll(c.Route().Path, "*", "")
			// Set doc url
			if len(cfg.URL) == 0 {
				cfg.URL = path.Join(prefix, defaultDocURL)
			}
		})

		var p string
		if p = utils.ImmutableString(c.Params("*")); p != "" {
			c.Path(p)
		} else {
			p = strings.TrimPrefix(c.Path(), prefix)
			p = strings.TrimPrefix(p, "/")
		}

		switch p {
		case defaultIndex:
			c.Type("html")
			return index.Execute(c, cfg)
		case defaultDocURL:
			doc, err := swag.ReadDoc()
			if err != nil {
				return err
			}
			return c.Type("json").SendString(doc)
		case "", "/":
			return c.Redirect(path.Join(prefix, defaultIndex), fiber.StatusMovedPermanently)
		default:
			return fs(c)
		}
	}
}
