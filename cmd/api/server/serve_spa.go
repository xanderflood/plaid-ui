package server

import (
	"github.com/gin-gonic/gin"
)

//ServeSPA serves the single-page app
func (a ServerAgent) ServeSPA(c *gin.Context) {
	a.renderer.RenderSPA(c)
}
