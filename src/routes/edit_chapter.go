package routes

import (
	"html"
	"strconv"
	"strings"

	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	"github.com/RyannKim327/constellation/src/backend"
	"github.com/RyannKim327/constellation/src/components"

	"github.com/gin-gonic/gin"
)

// EditChapter is a standalone page for rewriting an existing chapter.
var EditChapter = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/compose/chapter/:id/:chapter/edit",
	Handler:    editChapterPage,
	Permission: utils.PERMISSION_ALL,
	Title:      "Edit chapter · Istorya",
}

func editChapterPage(ctx *gin.Context) (string, utils.PARAMS) {
	chapterNumber, err := strconv.Atoi(ctx.Param("chapter"))
	novel, novelErr := backend.GetNovel(ctx.Param("id"))
	if err != nil || chapterNumber < 1 || novelErr != nil || chapterNumber > len(novel.Chapters) {
		return utils.HTParser(`
			<div class="container">
				{header}
				<section style="padding-block: 6rem;">
					<span class="eyebrow">The writing desk</span>
					<h1 class="display mt-4 mb-5" style="font-size: clamp(2.4rem, 5vw, 3.4rem);">Chapter not found.</h1>
					<p class="lead mb-8">This chapter may have been removed or the link is incorrect.</p>
					<a class="btn btn-primary" href="/read">Back to the library</a>
				</section>
			</div>
		`, utils.PARAMS{"header": components.Header(ctx.Request.URL.Path)}), nil
	}

	chapter := novel.Chapters[chapterNumber-1]
	chapterParam := strconv.Itoa(chapterNumber)

	return utils.HTParser(`
		<div class="container">
			{header}

			<header style="padding-block: 4rem 2.5rem;">
				<a class="back-link" href="/read/{id}?chapter={chapter}">← Back to {novel}</a>
				<span class="eyebrow mt-5" style="display:block;">The writing desk</span>
				<h1 class="display mt-4 mb-5" style="font-size: clamp(2.6rem, 5vw, 3.8rem);">
					Edit chapter {chapter}.
				</h1>
				<p class="lead" style="max-width: 40rem;">
					Rewrite this chapter in Markdown. The live preview updates as you type.
				</p>
			</header>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<div class="grid items-start gap-10 pt-12 lg:grid-cols-[minmax(0,1.6fr)_minmax(16rem,1fr)]">
				<section>
					<form class="novel-form" action="/api/novels/{id}/chapters/{chapter}/update" method="post">
						<div class="field">
							<label for="chapter_title">Chapter title *</label>
							<input id="chapter_title" name="chapter_title" type="text" required value="{title}" />
						</div>
						<div class="field">
							<label for="content">Chapter text *</label>
							<textarea id="content" name="content" rows="16" required>{content}</textarea>
							<span class="form-hint">Use {https://example.com/music.mp3} for music and {stop} to stop it.</span>
						</div>
						<div class="flex flex-wrap gap-3">
							<button class="btn btn-primary" type="submit">Save changes</button>
							<a class="btn btn-outline" href="/read/{id}?chapter={chapter}">Cancel</a>
						</div>
					</form>
				</section>

				<aside class="side-panel h-full">
					{preview}
				</aside>
			</div>

			{preview_script}

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/compose/chapter?novel={id}">Add another chapter</a>
			</footer>
		</div>
	`, utils.PARAMS{
		"header":         components.Header(ctx.Request.URL.Path),
		"id":             html.EscapeString(novel.ID),
		"novel":          html.EscapeString(novel.Title),
		"chapter":        chapterParam,
		"title":          html.EscapeString(chapter.Title),
		"content":        html.EscapeString(chapterMarkdown(chapter)),
		"preview":        components.PreviewPanel("content"),
		"preview_script": components.LivePreviewScript(),
	}), nil
}

func chapterMarkdown(chapter backend.Chapter) string {
	parts := make([]string, 0, len(chapter.Paragraphs))
	for _, paragraph := range chapter.Paragraphs {
		parts = append(parts, paragraph.Content)
	}
	return strings.Join(parts, "\n\n")
}
