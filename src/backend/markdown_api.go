package backend

import (
	"net/http"

	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	storyutils "github.com/RyannKim327/constellation/src/utils"

	"github.com/gin-gonic/gin"
)

// MarkdownAPI renders story Markdown for the live editor preview and the reader.
var MarkdownAPI = middleware.Route{
	Method:     utils.METHOD_POST,
	Path:       "api/markdown",
	Permission: utils.PERMISSION_ALL,
	Handler:    renderMarkdown,
}

type markdownInput struct {
	Text string `json:"text"`
}

func renderMarkdown(ctx *gin.Context) (string, utils.PARAMS) {
	var input markdownInput
	if wantsJSON(ctx) {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return "", nil
		}
	} else {
		input.Text = ctx.PostForm("text")
	}

	ctx.JSON(http.StatusOK, gin.H{"html": storyutils.RenderStory(input.Text)})
	return "", nil
}
