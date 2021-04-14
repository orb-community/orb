package v1

import (
	"github.com/gin-gonic/gin"
)

// we could make one of these per service or one for all of v1 depending on coupling points
type ExampleHandler struct {
	//Put needed services here
}

func RegisterExampleHandler(r *gin.RouterGroup) error {
	n := &ExampleHandler{}
	r.GET("/", n.listExampleHandler)
	r.GET("/:id", n.getExampleHandler)
	r.PUT("/", n.createExampleHandler)
	r.DELETE("/:id", n.deleteExampleHandler)
	return nil
}

type uriParams struct {
	Id int `uri:"id"`
}

func (n *ExampleHandler) listExampleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"example_networks": []string{},
	})
}

func (n *ExampleHandler) getExampleHandler(c *gin.Context) {
	var p uriParams
	if err := c.BindUri(&p); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, p)
}

func (n *ExampleHandler) createExampleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"example_networks": []string{},
	})
}

func (n *ExampleHandler) deleteExampleHandler(c *gin.Context) {
	var p uriParams
	if err := c.BindUri(&p); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, p)
}
