package site

import (
	"io/fs"

	constellation "constellation/src"
)

func Endpoints(assets fs.FS) []constellation.Endpoint {
	return []constellation.Endpoint{
		{Path: "/{$}", Title: "Clio · A home for your stories", Handler: constellation.HandlerSetup{Web: indexPage}},
		{Path: "/about", Title: "About · Clio", Handler: constellation.HandlerSetup{Web: aboutPage}},
		{Path: "/compose", Title: "Compose · Clio", Handler: constellation.HandlerSetup{Web: composePage}},
		{Path: "/compose/chapter", Title: "Add chapter · Clio", Handler: constellation.HandlerSetup{Web: addChapterPage}},
		{Path: "/compose/chapter/{id}/{chapter}/edit", Title: "Edit chapter · Clio", Handler: constellation.HandlerSetup{Api: editChapterEndpoint}},
		{Path: "/read", Title: "Library · Clio", Handler: constellation.HandlerSetup{Web: libraryPage}},
		{Path: "/read/{id}", Title: "Read · Clio", Handler: constellation.HandlerSetup{Api: readerEndpoint}},
		{Path: "/static/{file}", Method: constellation.GET, Handler: constellation.HandlerSetup{Api: staticFile(assets)}},
		{Path: "/api", Method: constellation.GET, Handler: constellation.HandlerSetup{Api: apiIndex}},
		{Path: "/api/novels", Method: constellation.GET, Handler: constellation.HandlerSetup{Api: apiListNovels}},
		{Path: "/api/novels", Method: constellation.POST, Handler: constellation.HandlerSetup{Api: apiCreateNovel}},
		{Path: "/api/novels/{id}", Method: constellation.GET, Handler: constellation.HandlerSetup{Api: apiGetNovel}},
		{Path: "/api/novels/{id}/chapters", Method: constellation.POST, Handler: constellation.HandlerSetup{Api: apiAddChapter}},
		{Path: "/api/novels/{id}/chapters/{chapter}", Method: constellation.PUT, Handler: constellation.HandlerSetup{Api: apiUpdateChapter}},
		{Path: "/api/novels/{id}/chapters/{chapter}/update", Method: constellation.POST, Handler: constellation.HandlerSetup{Api: apiUpdateChapter}},
		{Path: "/api/markdown", Method: constellation.POST, Handler: constellation.HandlerSetup{Api: apiMarkdown}},
		{Path: "/{path...}", Method: constellation.GET, Handler: constellation.HandlerSetup{Api: apiNotFound}},
	}
}
