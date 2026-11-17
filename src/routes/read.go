package routes

import (
	"errors"
	"html"
	"strconv"
	"strings"

	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"
	"github.com/RyannKim327/constellation/src/backend"
	"github.com/RyannKim327/constellation/src/components"
	storyutils "github.com/RyannKim327/constellation/src/utils"

	"github.com/gin-gonic/gin"
)

// ReadIndex lists every novel saved in the library.
var ReadIndex = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/read",
	Handler:    readIndex,
	Permission: utils.PERMISSION_ALL,
	Title:      "Library · Istorya",
}

// ReadNovel renders a single novel, one chapter at a time.
var ReadNovel = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "/read/:id",
	Handler:    readNovel,
	Permission: utils.PERMISSION_ALL,
	Title:      "Read · Istorya",
}

func readIndex(ctx *gin.Context) (string, utils.PARAMS) {
	novels, err := backend.ListNovels()
	if err != nil {
		novels = nil
	}

	cards := strings.Builder{}
	for _, novel := range novels {
		description := novel.Description
		if description == "" {
			description = "A story waiting to be read."
		}
		cards.WriteString(utils.HTParser(`
			<a class="card library-card" href="/read/{id}">
				<span class="eyebrow">{chapters} chapters</span>
				<h2 class="card-title mt-3">{title}</h2>
				<p class="card-text">{description}</p>
				<span class="library-byline">by {author}</span>
			</a>
		`, utils.PARAMS{
			"id":          html.EscapeString(novel.ID),
			"chapters":    strconv.Itoa(len(novel.Chapters)),
			"title":       html.EscapeString(novel.Title),
			"description": html.EscapeString(description),
			"author":      html.EscapeString(defaultAuthor(novel.Author)),
		}))
	}

	content := cards.String()
	if content == "" {
		content = `<p class="lead">No stories yet. <a href="/compose">Write the first one</a>.</p>`
	}

	return utils.HTParser(`
		<div class="container">
			{header}

			<header style="padding-block: 4rem 2.5rem;">
				<span class="eyebrow">The library</span>
				<h1 class="display mt-4 mb-5" style="font-size: clamp(2.6rem, 5vw, 3.8rem);">
					Stories, ready to open.
				</h1>
				<p class="lead" style="max-width: 40rem;">
					Every novel you compose lives here. Choose a book and settle in for a chapter.
				</p>
			</header>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<section class="grid gap-6 py-12 md:grid-cols-2">
				{content}
			</section>

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/compose">Compose a novel</a>
			</footer>
		</div>
	`, utils.PARAMS{
		"header":  components.Header(ctx.Request.URL.Path),
		"content": content,
	}), nil
}

func readNovel(ctx *gin.Context) (string, utils.PARAMS) {
	novel, err := backend.GetNovel(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, backend.ErrNovelNotFound) {
			return utils.HTParser(`
				<div class="container">
					{header}
					<section style="padding-block: 6rem;">
						<span class="eyebrow">Not in the library</span>
						<h1 class="display mt-4 mb-5" style="font-size: clamp(2.4rem, 5vw, 3.4rem);">
							This story could not be found.
						</h1>
						<p class="lead mb-8">It may have been moved, or the link may be incomplete.</p>
						<a class="btn btn-primary" href="/read">Back to the library</a>
					</section>
				</div>
			`, utils.PARAMS{"header": components.Header(ctx.Request.URL.Path)}), nil
		}
		return "", nil
	}

	if len(novel.Chapters) == 0 {
		return "", nil
	}

	chapterNumber, err := strconv.Atoi(ctx.Query("chapter"))
	if err != nil || chapterNumber < 1 {
		chapterNumber = 1
	}
	if chapterNumber > len(novel.Chapters) {
		chapterNumber = len(novel.Chapters)
	}
	chapter := novel.Chapters[chapterNumber-1]

	paragraphs := strings.Builder{}
	for _, paragraph := range chapter.Paragraphs {
		paragraphs.WriteString(storyutils.RenderStory(paragraph.Content))
	}

	toc := strings.Builder{}
	for i, item := range novel.Chapters {
		active := ""
		if i+1 == chapterNumber {
			active = " active"
		}
		toc.WriteString(utils.HTParser(`
			<li><a class="toc-link{active}" href="/read/{id}?chapter={n}">{title}</a></li>
		`, utils.PARAMS{
			"active": active,
			"id":     html.EscapeString(novel.ID),
			"n":      strconv.Itoa(i + 1),
			"title":  html.EscapeString(item.Title),
		}))
	}

	prev := ""
	if chapterNumber > 1 {
		prev = utils.HTParser(`<a class="btn btn-outline" href="/read/{id}?chapter={n}">← Previous chapter</a>`, utils.PARAMS{
			"id": html.EscapeString(novel.ID), "n": strconv.Itoa(chapterNumber - 1),
		})
	}
	next := ""
	if chapterNumber < len(novel.Chapters) {
		next = utils.HTParser(`<a class="btn btn-primary" href="/read/{id}?chapter={n}">Next chapter →</a>`, utils.PARAMS{
			"id": html.EscapeString(novel.ID), "n": strconv.Itoa(chapterNumber + 1),
		})
	}

	description := novel.Description
	if description == "" {
		description = "A story from the Istorya library."
	}

	return utils.HTParser(`
		<div class="container">
			{header}

			<div class="start-overlay" id="start-overlay">
				<div class="start-card">
					<span class="eyebrow">Before you begin</span>
					<h2 class="display mt-4 mb-3">{title}</h2>
					<p class="muted mb-7">Start reading to let the story play its music.</p>
					<button class="btn btn-primary" id="start-reading" type="button">Start reading</button>
				</div>
			</div>

			<audio id="story-audio" loop preload="none"></audio>
			<div class="now-playing" id="now-playing" hidden>
				<span class="now-playing-label">♪ Now playing</span>
				<span class="now-playing-title" id="now-title"></span>
				<button class="player-btn" id="audio-toggle" type="button" aria-label="Pause">❚❚</button>
				<button class="player-btn" id="audio-mute" type="button" aria-label="Mute">Mute</button>
			</div>

			<header style="padding-block: 3rem 2rem;">
				<a class="back-link" href="/read">← Back to the library</a>
				<h1 class="display mt-5 mb-3" style="font-size: clamp(2.4rem, 5vw, 3.6rem);">{title}</h1>
				<p class="muted" style="font-family: var(--font-ui); font-size: 0.85rem;">by {author}</p>
				<p class="lead mt-4" style="max-width: 42rem;">{description}</p>
			</header>

			<div class="divider"><span class="divider-mark">✦</span></div>

			<div class="grid items-start gap-14 py-12 lg:grid-cols-[minmax(0,2fr)_minmax(15rem,1fr)]">
				<article class="">
					<div class="reader-chapter-head sticky top-0">
						<span class="eyebrow">Chapter {number} of {total}</span>
						<h2 class="display mt-3 mb-6" style="font-size: 2rem;">{chapter_title}</h2>
					</div>
					<div class="reader-body">{paragraphs}</div>
					<div class="reader-nav">{prev}{next}</div>
				</article>

				<aside class="side-panel">
					<p class="eyebrow" style="margin-bottom: 1rem;">Contents</p>
					<ol class="toc">{toc}</ol>
					<div class="flex flex-col gap-2 mt-6">
						<a class="btn btn-primary" href="/compose/chapter?novel={id}">Add a chapter</a>
						<a class="btn btn-outline" href="/compose/chapter/{id}/{number}/edit">Edit this chapter</a>
					</div>
				</aside>
			</div>

			<footer class="site-footer">
				<span>Istorya · A home for your stories</span>
				<a href="/read">The library</a>
			</footer>

			<script>
			(function () {
				var overlay = document.getElementById('start-overlay');
				var startBtn = document.getElementById('start-reading');
				var audio = document.getElementById('story-audio');
				var bar = document.getElementById('now-playing');
				var title = document.getElementById('now-title');
				var toggle = document.getElementById('audio-toggle');
				var mute = document.getElementById('audio-mute');
				var markers = Array.prototype.slice.call(document.querySelectorAll('[data-audio-src]'));
				var stops = Array.prototype.slice.call(document.querySelectorAll('[data-audio-stop]'));
				var unlocked = false;
				var activeSrc = null;
				var resumeOnVisible = false;

				function trackName(src) {
					return src.split('?')[0].split('/').pop() || 'Story music';
				}
				function updateBar(playing) {
					if (!activeSrc) { bar.hidden = true; return; }
					bar.hidden = false;
					title.textContent = trackName(activeSrc);
					toggle.textContent = playing ? '❚❚' : '▶';
				}
				function playSrc(src) {
					if (!unlocked || !src) return;
					if (activeSrc !== src) {
						activeSrc = src;
						audio.src = src;
						audio.loop = true;
						audio.load();
					}
					audio.play().then(function () { updateBar(true); }).catch(function () { updateBar(false); });
				}
				function stopAudio() {
					audio.pause();
					audio.currentTime = 0;
					activeSrc = null;
					updateBar(false);
				}
				function currentMarker() {
					for (var i = 0; i < markers.length; i++) {
						var box = markers[i].getBoundingClientRect();
						if (box.top < window.innerHeight && box.bottom > 0) return markers[i];
					}
					return null;
				}

				startBtn.addEventListener('click', function () {
					unlocked = true;
					overlay.hidden = true;
					var marker = currentMarker();
					if (marker) playSrc(marker.getAttribute('data-audio-src'));
				});
				toggle.addEventListener('click', function () {
					if (!activeSrc) return;
					if (audio.paused) {
						audio.play().then(function () { updateBar(true); }).catch(function () {});
					} else {
						audio.pause();
						updateBar(false);
					}
				});
				mute.addEventListener('click', function () {
					audio.muted = !audio.muted;
					mute.textContent = audio.muted ? 'Unmute' : 'Mute';
				});
				audio.addEventListener('play', function () { updateBar(true); });
				audio.addEventListener('pause', function () { updateBar(false); });

				if ('IntersectionObserver' in window) {
					var observer = new IntersectionObserver(function (entries) {
						if (!unlocked) return;
						entries.forEach(function (entry) {
							if (!entry.isIntersecting) return;
							if (entry.target.hasAttribute('data-audio-stop')) { stopAudio(); return; }
							var src = entry.target.getAttribute('data-audio-src');
							if (src) playSrc(src);
						});
					}, { threshold: 0.5 });
					markers.concat(stops).forEach(function (el) { observer.observe(el); });
				}

				document.addEventListener('visibilitychange', function () {
					if (!unlocked) return;
					if (document.hidden) {
						if (!audio.paused) { resumeOnVisible = true; audio.pause(); }
					} else if (resumeOnVisible && activeSrc) {
						resumeOnVisible = false;
						audio.play().catch(function () {});
					}
				});
			})();
			</script>
		</div>
	`, utils.PARAMS{
		"header":        components.Header(ctx.Request.URL.Path),
		"id":            html.EscapeString(novel.ID),
		"title":         html.EscapeString(novel.Title),
		"author":        html.EscapeString(defaultAuthor(novel.Author)),
		"description":   html.EscapeString(description),
		"number":        strconv.Itoa(chapterNumber),
		"total":         strconv.Itoa(len(novel.Chapters)),
		"chapter_title": html.EscapeString(chapter.Title),
		"paragraphs":    paragraphs.String(),
		"prev":          prev,
		"next":          next,
		"toc":           toc.String(),
	}), nil
}

func defaultAuthor(author string) string {
	if author == "" {
		return "Unknown author"
	}
	return author
}
