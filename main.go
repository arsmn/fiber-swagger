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

// New returns custom handler
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

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
