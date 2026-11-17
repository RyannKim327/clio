package routes

import (
	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	"github.com/RyannKim327/constellation/src/components"

	"github.com/gin-gonic/gin"
)

// Compose is the writing desk for starting a new novel.
var Compose = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/compose",
	Handler:    compose,
	Permission: utils.PERMISSION_ALL,
	Title:      "Compose · Istorya",
}

func compose(ctx *gin.Context) (string, utils.PARAMS) {
	return utils.HTParser(`
		<div class="container">
			{header}

			<header style="padding-block: 4rem 2.5rem;">
				<span class="eyebrow">The writing desk</span>
				<h1 class="display mt-4 mb-5" style="font-size: clamp(2.6rem, 5vw, 3.8rem);">
					Begin with a blank page.
				</h1>
				<p class="lead" style="max-width: 40rem;">
					Give your story a title and a first chapter. Everything you write is saved
					to the library, ready to read whenever you return.
				</p>
			</header>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<div class="grid items-start gap-10 pt-12 lg:grid-cols-[minmax(0,1.6fr)_minmax(16rem,1fr)]">
				<section>
					<h2 class="display mb-2" style="font-size: 1.8rem;">Start a new novel</h2>
					<p class="muted mb-7" style="font-size: 0.92rem;">Fields marked with * are required.</p>

					<form class="novel-form" action="/api/novels" method="post">
						<div class="field">
							<label for="title">Novel title *</label>
							<input id="title" name="title" type="text" required placeholder="The name of your story" />
						</div>
						<div class="field-grid">
							<div class="field">
								<label for="author">Author</label>
								<input id="author" name="author" type="text" placeholder="Your name" />
							</div>
							<div class="field">
								<label for="chapter_title">Chapter title</label>
								<input id="chapter_title" name="chapter_title" type="text" placeholder="Chapter 1" />
							</div>
						</div>
						<div class="field">
							<label for="description">A short description</label>
							<textarea id="description" name="description" rows="2" placeholder="What is this story about?"></textarea>
						</div>
						<div class="field">
							<label for="content">First chapter *</label>
							<textarea id="content" name="content" rows="10" required placeholder="Write your opening here..."></textarea>
							<span class="form-hint">Markdown is live on the right. Use **bold**, *italic*, # headings, and {https://example.com/music.mp3} for music.</span>
						</div>
						<button class="btn btn-primary" type="submit">Save novel</button>
					</form>
				</section>

				<aside class="side-panel h-full">
					{preview}
				</aside>
			</div>

			{preview_script}

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/read">Read the library</a>
			</footer>
		</div>
	`, utils.PARAMS{
		"header":         components.Header(ctx.Request.URL.Path),
		"preview":        components.PreviewPanel("content"),
		"preview_script": components.LivePreviewScript(),
	}), nil
}
