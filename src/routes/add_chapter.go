package routes

import (
	"html"
	"strings"

	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	"github.com/RyannKim327/constellation/src/backend"
	"github.com/RyannKim327/constellation/src/components"

	"github.com/gin-gonic/gin"
)

// AddChapter is a standalone page for appending a chapter to a saved novel.
var AddChapter = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/compose/chapter",
	Handler:    addChapterPage,
	Permission: utils.PERMISSION_ALL,
	Title:      "Add chapter · Istorya",
}

func addChapterPage(ctx *gin.Context) (string, utils.PARAMS) {
	novels, err := backend.ListNovels()
	if err != nil {
		novels = nil
	}

	selected := ctx.Query("novel")
	options := strings.Builder{}
	options.WriteString(`<option value="" disabled selected>Choose a novel</option>`)
	firstID := ""
	for _, novel := range novels {
		if firstID == "" {
			firstID = novel.ID
		}
		chosen := ""
		if novel.ID == selected {
			chosen = " selected"
		}
		options.WriteString(utils.HTParser(`
			<option value="{id}"{chosen}>{title}</option>
		`, utils.PARAMS{
			"id":     html.EscapeString(novel.ID),
			"chosen": chosen,
			"title":  html.EscapeString(novel.Title),
		}))
	}

	if selected == "" {
		selected = firstID
	}

	if len(novels) == 0 {
		return utils.HTParser(`
			<div class="container">
				{header}
				<section style="padding-block: 6rem;">
					<span class="eyebrow">The writing desk</span>
					<h1 class="display mt-4 mb-5" style="font-size: clamp(2.4rem, 5vw, 3.4rem);">No stories yet.</h1>
					<p class="lead mb-8">Start a novel before adding chapters to it.</p>
					<a class="btn btn-primary" href="/compose">Compose a novel</a>
				</section>
			</div>
		`, utils.PARAMS{"header": components.Header(ctx.Request.URL.Path)}), nil
	}

	return utils.HTParser(`
		<div class="container">
			{header}

			<header style="padding-block: 4rem 2.5rem;">
				<a class="back-link" href="/read">← Back to the library</a>
				<span class="eyebrow mt-5" style="display:block;">The writing desk</span>
				<h1 class="display mt-4 mb-5" style="font-size: clamp(2.6rem, 5vw, 3.8rem);">
					Add a chapter.
				</h1>
				<p class="lead" style="max-width: 40rem;">
					Continue a story already in your library. Write in Markdown and watch the
					page take shape as you type.
				</p>
			</header>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<div class="grid items-start gap-10 pt-12 lg:grid-cols-[minmax(0,1.6fr)_minmax(16rem,1fr)]">
				<section>
					<form class="novel-form" action="/api/novels/{selected}/chapters" method="post" id="chapter-form">
						<div class="field">
							<label for="novel-select">Novel *</label>
							<select id="novel-select" name="novel_id" required>{options}</select>
						</div>
						<div class="field">
							<label for="chapter_title">Chapter title</label>
							<input id="chapter_title" name="chapter_title" type="text" placeholder="Chapter 2" />
							<span class="form-hint">Leave blank to number it automatically.</span>
						</div>
						<div class="field">
							<label for="content">Chapter text *</label>
							<textarea id="content" name="content" rows="14" required placeholder="Continue the story..."></textarea>
							<span class="form-hint">Use {https://example.com/music.mp3} for music and {stop} to stop it.</span>
						</div>
						<div class="flex flex-wrap gap-3">
							<button class="btn btn-primary" type="submit">Save chapter</button>
							<a class="btn btn-outline" href="/read">Cancel</a>
						</div>
					</form>
					<script>
						var select = document.getElementById('novel-select');
						var form = document.getElementById('chapter-form');
						select.addEventListener('change', function () {
							form.action = '/api/novels/' + this.value + '/chapters';
						});
					</script>
				</section>

				<aside class="side-panel h-full">
					{preview}
				</aside>
			</div>

			{preview_script}

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/compose">Start a new novel</a>
			</footer>
		</div>
	`, utils.PARAMS{
		"header":         components.Header(ctx.Request.URL.Path),
		"options":        options.String(),
		"selected":       html.EscapeString(selected),
		"preview":        components.PreviewPanel("content"),
		"preview_script": components.LivePreviewScript(),
	}), nil
}
