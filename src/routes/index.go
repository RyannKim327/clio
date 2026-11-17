package routes

import (
	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	"github.com/RyannKim327/constellation/src/components"

	"github.com/gin-gonic/gin"
)

// Index is the landing page. It introduces Istorya as a home for drafting and
// designing novels.
var Index = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/",
	Handler:    layout,
	Permission: utils.PERMISSION_ALL,
	Title:      "Istorya · A home for your stories",
}

func layout(ctx *gin.Context) (string, utils.PARAMS) {
	return utils.HTParser(`
		<div class="container">
			{header}

			<section class="grid items-center gap-14 pt-20 pb-16 lg:grid-cols-[1.08fr_0.92fr]">
				<div>
					<span class="eyebrow">A writing studio for storytellers</span>
					<h1 class="display mt-5 mb-6" style="font-size: clamp(2.8rem, 6vw, 4.4rem);">
						Stories deserve a place to <em style="color: var(--accent); font-style: italic;">grow.</em>
					</h1>
					<p class="lead mb-8" style="max-width: 34rem;">
						Istorya is a web-based studio for designing novels. Shape your chapters,
						build your world, and watch a simple idea become a book — one beautiful
						page at a time.
					</p>
					<div class="flex flex-wrap items-center gap-3">
						<a class="btn btn-primary" href="#features">Begin a story</a>
						<a class="btn btn-outline" href="/about">Our story</a>
					</div>
					<p class="muted mt-6" style="font-family: var(--font-ui); font-size: 0.8rem;">
						Made for tales of every length and genre.
					</p>
				</div>

				<div class="manuscript">
					<span class="manuscript-kicker">Chapter One</span>
					<h2 class="manuscript-title">The Quiet Between Words</h2>
					<div class="manuscript-rule"></div>
					<p class="manuscript-text">
						There is a moment before a story begins — when the page is still
						empty, and everything is still possible. Istorya begins there, with you.
					</p>
					<div class="manuscript-meta">
						<span>Draft · 1,284 words</span>
						<span>Istorya</span>
					</div>
				</div>
			</section>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<section id="features" class="py-16">
				<span class="eyebrow">Everything in its place</span>
				<h2 class="display mt-4 mb-3" style="font-size: clamp(2rem, 4vw, 2.8rem);">
					From first draft to finished book
				</h2>
				<p class="lead mb-10" style="max-width: 40rem;">
					The tools you need to write, organize, and present a novel — without
					getting in the way of the story itself.
				</p>

				<div class="grid gap-6 md:grid-cols-3">
					<article class="card">
						<span class="eyebrow">01 / Outline</span>
						<h3 class="card-title mt-3">Build the frame</h3>
						<p class="card-text">
							Map chapters, scenes, and characters in a structure that bends to
							the way you write, not the other way around.
						</p>
					</article>
					<article class="card">
						<span class="eyebrow">02 / Design</span>
						<h3 class="card-title mt-3">Dress the page</h3>
						<p class="card-text">
							Choose refined typography and gentle layouts that give your words
							the room and the feeling they deserve.
						</p>
					</article>
					<article class="card">
						<span class="eyebrow">03 / Share</span>
						<h3 class="card-title mt-3">Open the book</h3>
						<p class="card-text">
							Turn your manuscript into a reading experience you can share with
							beta readers, friends, or the world.
						</p>
					</article>
				</div>
			</section>

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/about">About the project</a>
			</footer>
		</div>
	`, utils.PARAMS{
		"header": components.Header(ctx.Request.URL.Path),
	}), nil
}
