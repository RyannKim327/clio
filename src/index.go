package src

import (
	routes "github.com/RyannKim327/constellation/src/routes"

	c "github.com/RyannKim327/constellation/config"
	mw "github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/src/backend"
)

// TODO: Collect all the Routes here
// NOTE: I suggest to use routes (folder) to make it organize
var _routes = []mw.Route{
	backend.API,
	backend.NovelsAPI,
	backend.NovelAPI,
	backend.CreateNovelAPI,
	backend.AddChapterAPI,
	backend.UpdateChapterAPI,
	backend.UpdateChapterFormAPI,
	backend.MarkdownAPI,

	routes.Index,
	routes.About,
	routes.Compose,
	routes.AddChapter,
	routes.EditChapter,
	routes.ReadIndex,
	routes.ReadNovel,
}

func Main() {
	// TODO: This code will going to initiate everything
	c.Constellation(StaticFS, _routes)
}
