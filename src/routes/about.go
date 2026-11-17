package routes

import (
	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	"github.com/RyannKim327/constellation/src/components"

	"github.com/gin-gonic/gin"
)

// About describes Istorya, its purpose, and what a writer can do with it.
var About = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/about",
	Handler:    about,
	Permission: utils.PERMISSION_ALL,
	Title:      "About · Istorya",
}

func about(ctx *gin.Context) (string, utils.PARAMS) {
	return utils.HTParser(`
		<div class="container">
			{header}

			<header style="padding-block: 4rem 2.5rem;">
				<span class="eyebrow">About the project</span>
				<h1 class="display mt-4 mb-5" style="font-size: clamp(2.6rem, 5vw, 3.8rem); max-width: 20ch;">
					Stories, given room to breathe.
				</h1>
				<p class="lead" style="max-width: 42rem;">
					Istorya is a web-based app for designing novels — a quiet place to plan,
					write, and shape a book from the very first idea to the finished page.
				</p>
			</header>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<div class="grid items-start gap-14 pt-12 lg:grid-cols-[minmax(0,2fr)_minmax(16rem,1fr)]">
				<section class="prose">
					<h2 style="margin-top: 0;">Why it exists</h2>
					<p>
						Every novel starts as something small: a character, a place, a single
						line you can't shake. Istorya was made to protect that spark. Instead of
						fighting with formatting and scattered notes, you get a calm space where
						the writing stays at the center.
					</p>
					<p>
						The name comes from the Filipino word for <em>story</em>, and the idea is
						simple — give storytellers the tools to build something lasting, without
						ever getting between them and the words.
					</p>

					<h2>What you can do</h2>
					<ul>
						<li>Outline chapters and scenes at your own pace.</li>
						<li>Keep characters, settings, and notes close to the draft.</li>
						<li>Preview your novel with book-like typography and layout.</li>
						<li>Refine the reading experience before you share it.</li>
					</ul>

					<h2>How it feels</h2>
					<p>
						Pages are warm like paper, type is set in a timeless serif, and the
						interface stays out of the way. It should feel less like using software
						and more like opening a notebook you're happy to return to.
					</p>

					<div class="callout">
						<p class="eyebrow" style="margin-bottom: 0.4rem;">Begin anywhere</p>
						<p>
							You don't need a complete plan. Bring a fragment, a world, or just a
							title — Istorya gives it somewhere to live.
						</p>
					</div>
				</section>

				<aside class="side-panel">
					<p class="eyebrow" style="margin-bottom: 1rem;">At a glance</p>
					<dl style="margin: 0 0 1.5rem;">
						<div class="stat-row"><dt>Made for</dt><dd>Novelists</dd></div>
						<div class="stat-row"><dt>Workspace</dt><dd>Web based</dd></div>
						<div class="stat-row"><dt>Focus</dt><dd>Draft &amp; design</dd></div>
						<div class="stat-row"><dt>Typography</dt><dd>Classic serif</dd></div>
						<div class="stat-row"><dt>Style</dt><dd>Warm &amp; quiet</dd></div>
					</dl>
					<a class="btn btn-primary" href="/">Back to home</a>
				</aside>
			</div>

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/">Home</a>
			</footer>
		</div>
	`, utils.PARAMS{
		"header": components.Header(ctx.Request.URL.Path),
	}), nil
}
