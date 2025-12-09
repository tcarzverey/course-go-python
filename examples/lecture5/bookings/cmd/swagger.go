package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tcarzverey/bookings/internal/generated/api"
)

const swaggerUIPage = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: "/swagger.json",
          dom_id: "#swagger-ui",
        });
      };
    </script>
  </body>
</html>`

func initSwagger(r *gin.Engine) {
	// Swagger JSON
	r.GET("/swagger.json", func(c *gin.Context) {
		swagger, err := api.GetSwagger()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get swagger spec",
			})
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, swagger)
	})

	// Swagger UI HTML
	r.GET("/swagger", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerUIPage)
	})

	// Альтернативно: редирект на /swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger")
	})
}
