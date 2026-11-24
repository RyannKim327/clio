package site

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	constellation "constellation/src"
	"constellation/src/elements"
)

func indexPage(r *constellation.Request) string {
	return page(r,
		section("py-14 sm:py-20",
			div("grid items-center gap-14 lg:grid-cols-[1.05fr_0.95fr]",
				div("space-y-7",
					eyebrow("A writing studio for storytellers"),
					elements.H1(constellation.WebComponent{
						Class: "display text-5xl leading-[1.08] sm:text-6xl",
						Children: []any{
							"Stories deserve a place to ",
							elements.Em(constellation.WebComponent{Class: "accent-italic", Text: "grow."}),
						},
					}),
					lead("Clio is a quiet corner of the web for drafting novels. Write in Markdown, watch the page come alive, and keep every chapter in a library of its own."),
					div("flex flex-wrap gap-3 pt-2",
						buttonLink("#features", "btn-primary", "Begin a story"),
						buttonLink("/about", "btn-outline", "Our story"),
					),
					paragraph("text-sm text-slate-500", "Made for tales of every length and genre."),
				),
				manuscriptCard(),
			),
		),
		divider(),
		featuresSection(),
	)
}

func manuscriptCard() string {
	return div("manuscript-shell",
		div("manuscript",
			span("manuscript-kicker", "Chapter One"),
			heading("h2", "manuscript-title", "The Quiet Between Words"),
			div("manuscript-rule"),
			paragraph("manuscript-text", "There was a particular kind of silence in the room before the first word arrived — the kind that felt less like waiting and more like listening."),
			div("manuscript-meta",
				span("", "Draft · 1,284 words"),
				span("", "Clio"),
			),
		),
	)
}

func featuresSection() string {
	features := []struct {
		number string
		title  string
		text   string
	}{
		{number: "01 / Outline", title: "Build the frame", text: "Give a story its shape. Add chapters, reorder your thoughts, and let the draft find its own rhythm."},
		{number: "02 / Design", title: "Dress the page", text: "Write in Markdown with a live preview built for long-form reading — headings, quotes, code, and song cues."},
		{number: "03 / Share", title: "Open the book", text: "A reader with a table of contents, chapter navigation, and optional music that follows the paragraphs."},
	}

	cards := make([]any, 0, len(features))
	for _, feature := range features {
		cards = append(cards,
			div("feature-card",
				span("feature-number", feature.number),
				heading("h3", "card-title", feature.title),
				paragraph("card-text", feature.text),
			),
		)
	}

	return elements.Section(constellation.WebComponent{
		Class:      "py-6",
		Attributes: map[string]string{"id": "features"},
		Children: []any{
			div("flex flex-wrap items-end justify-between gap-4",
				div("space-y-3",
					eyebrow("The studio"),
					heading("h2", "display text-3xl sm:text-4xl", "From first draft to finished book"),
				),
				buttonLink("/compose", "btn-outline", "Open the writing desk"),
			),
			div("mt-10 grid gap-5 md:grid-cols-3", cards...),
		},
	})
}

func aboutPage(r *constellation.Request) string {
	return page(r,
		section("py-14",
			eyebrow("About the project"),
			heading("h1", "display mt-5 text-5xl", "Stories, given room to breathe."),
			lead("Clio is a small writing studio for people who still believe a paragraph can hold a whole world."),
			divider(),
			div("grid gap-12 lg:grid-cols-[1.2fr_0.8fr]",
				div("prose space-y-8",
					div("space-y-3",
						heading("h2", "prose-title", "Why it exists"),
						paragraph("", "The name is borrowed from the Filipino word for story. This project began as a place to keep drafts close at hand — simple enough to open on a whim, careful enough to make the words feel at home."),
					),
					div("space-y-3",
						heading("h2", "prose-title", "What you can do"),
						elements.Ul(constellation.WebComponent{Class: "prose-list", Children: []any{
							elements.Li(constellation.WebComponent{Text: "Outline chapters and draft them one at a time."}),
							elements.Li(constellation.WebComponent{Text: "Keep notes, characters, and settings in the same manuscript."}),
							elements.Li(constellation.WebComponent{Text: "Preview Markdown with book-like typography as you type."}),
							elements.Li(constellation.WebComponent{Text: "Refine a chapter before opening it in the reader."}),
						}}),
					),
					div("space-y-3",
						heading("h2", "prose-title", "How it feels"),
						paragraph("", "Warm like paper, steady like a desk lamp, and quiet enough to forget the interface. A serif page, soft margins, and just enough ornament to mark the way."),
					),
					div("callout",
						heading("h3", "callout-title", "Begin anywhere"),
						paragraph("", "Bring a fragment, a world, or just a title. The first chapter can always change its name."),
					),
				),
				panel("h-fit p-7",
					span("eyebrow", "At a glance"),
					elements.Dl(constellation.WebComponent{Class: "stat-list mt-6", Children: []any{
						statRow("Made for", "Novelists"),
						statRow("Workspace", "Web based"),
						statRow("Focus", "Draft & design"),
						statRow("Typography", "Classic serif"),
						statRow("Style", "Warm & quiet"),
					}}),
					div("mt-8", buttonLink("/", "btn-outline", "Back to home")),
				),
			),
		),
	)
}

func statRow(term, description string) string {
	return div("stat-row",
		elements.Dt(constellation.WebComponent{Text: term}),
		elements.Dd(constellation.WebComponent{Text: description}),
	)
}

func composePage(r *constellation.Request) string {
	return page(r,
		section("py-14",
			eyebrow("The writing desk"),
			heading("h1", "display mt-5 text-5xl", "Begin with a blank page."),
			lead("Name the story, set the scene, and write the opening. Everything you save lands in the library, ready to be read chapter by chapter."),
			divider(),
			div("grid items-start gap-10 lg:grid-cols-[1.05fr_0.95fr]",
				panel("p-6 sm:p-8",
					heading("h2", "panel-title", "Start a new novel"),
					paragraph("form-note", "Fields marked with * are required."),
					elements.Form(constellation.WebComponent{
						Class: "novel-form mt-7",
						Attributes: map[string]string{
							"action": "/api/novels",
							"method": "post",
						},
						Children: []any{
							textField("title", "Novel title *", "title", "The name of your story", "", "", true),
							div("field-grid",
								textField("author", "Author", "author", "Your name", "", "", false),
								textField("chapter_title", "Chapter title", "chapter_title", "Chapter 1", "", "", false),
							),
							textareaField("description", "Description", "description", "What is this story about?", "", "", 2, false),
							textareaField("content", "Opening chapter *", "content", "Write your opening here...", "", "Markdown is live on the right. Try **bold**, *italic*, # headings, and {https://example.com/music.mp3} for music.", 10, true),
							elements.Button(constellation.WebComponent{
								Class:      "btn btn-primary w-fit",
								Attributes: map[string]string{"type": "submit"},
								Text:       "Save novel",
							}),
						},
					}),
				),
				div("space-y-5",
					previewPanel("content"),
					panel("tip-card p-5", span("eyebrow", "A small hint"), paragraph("mt-3 text-sm", "Use {stop} anywhere in a chapter to fade the music out when the scene turns.")),
				),
			),
			div("mt-10", buttonLink("/read", "text-link", "Read the library →")),
		),
		previewScript(),
	)
}

func addChapterPage(r *constellation.Request) string {
	novels, err := ListNovels()
	if err != nil || len(novels) == 0 {
		return page(r,
			section("py-24 text-center",
				eyebrow("The writing desk"),
				heading("h1", "display mt-5 text-4xl", "No stories yet."),
				paragraph("lead mx-auto mt-4", "Create a novel first, then come back to continue it."),
				div("mt-8", buttonLink("/compose", "btn-primary", "Compose a novel")),
			),
		)
	}

	selected := r.URL.Query().Get("novel")
	found := false
	for _, novel := range novels {
		if novel.ID == selected {
			found = true
			break
		}
	}
	if !found {
		selected = novels[0].ID
	}

	options := make([]any, 0, len(novels)+1)
	placeholder := map[string]string{"value": "", "disabled": "disabled"}
	if r.URL.Query().Get("novel") == "" {
		placeholder["selected"] = "selected"
	}
	options = append(options, elements.Option(constellation.WebComponent{
		Text:       "Choose a novel",
		Attributes: placeholder,
	}))
	for _, novel := range novels {
		attributes := map[string]string{"value": novel.ID}
		if novel.ID == selected {
			attributes["selected"] = "selected"
		}
		options = append(options, elements.Option(constellation.WebComponent{
			Text:       novel.Title,
			Attributes: attributes,
		}))
	}

	return page(r,
		section("py-14",
			buttonLink("/read", "back-link", "← Back to the library"),
			div("mt-10 space-y-5",
				eyebrow("The writing desk"),
				heading("h1", "display text-5xl", "Add a chapter."),
				lead("Pick up where you left off. The preview follows every word as you write."),
			),
			divider(),
			div("grid items-start gap-10 lg:grid-cols-[1.05fr_0.95fr]",
				panel("p-6 sm:p-8",
					elements.Form(constellation.WebComponent{
						Class: "novel-form",
						Attributes: map[string]string{
							"id":     "chapter-form",
							"action": "/api/novels/" + url.PathEscape(selected) + "/chapters",
							"method": "post",
						},
						Children: []any{
							selectField("novel-select", "Novel", "novel_id", options...),
							textField("chapter_title", "Chapter title", "chapter_title", "Chapter 2", "", "Leave blank to number it automatically.", false),
							textareaField("content", "Chapter *", "content", "Continue the story...", "", "Markdown works here too. Use {https://example.com/music.mp3} and {stop} for music cues.", 14, true),
							div("form-actions",
								elements.Button(constellation.WebComponent{
									Class:      "btn btn-primary",
									Attributes: map[string]string{"type": "submit"},
									Text:       "Save chapter",
								}),
								buttonLink("/read", "btn btn-outline", "Cancel"),
							),
						},
					}),
				),
				previewPanel("content"),
			),
			div("mt-10", buttonLink("/compose", "text-link", "Start a new novel →")),
		),
		chapterSwitchScript(),
		previewScript(),
	)
}

func editChapterPage(r *constellation.Request) (string, int) {
	id := pathValue(r, "id")
	number, parseErr := strconv.Atoi(pathValue(r, "chapter"))
	novel, novelErr := GetNovel(id)

	if parseErr != nil || novelErr != nil || number < 1 || number > len(novel.Chapters) {
		return notFoundPage(r, "Chapter not found.", "/read", "Back to the library"), http.StatusNotFound
	}

	chapter := novel.Chapters[number-1]
	storyPath := "/read/" + url.PathEscape(novel.ID)
	chapterURL := storyPath + "?chapter=" + strconv.Itoa(number)
	chapterTitle := strings.TrimSpace(chapter.Title)
	if chapterTitle == "" {
		chapterTitle = "Chapter " + strconv.Itoa(number)
	}

	return page(r,
		section("py-14",
			buttonLink(chapterURL, "back-link", "← Back to "+novel.Title),
			div("mt-10 space-y-5",
				eyebrow("The writing desk"),
				heading("h1", "display text-5xl", "Edit chapter "+strconv.Itoa(number)+"."),
				lead("Refine the words, then open the reader to hear how they land."),
			),
			divider(),
			div("grid items-start gap-10 lg:grid-cols-[1.05fr_0.95fr]",
				panel("p-6 sm:p-8",
					elements.Form(constellation.WebComponent{
						Class: "novel-form",
						Attributes: map[string]string{
							"action": "/api/novels/" + url.PathEscape(novel.ID) + "/chapters/" + strconv.Itoa(number) + "/update",
							"method": "post",
						},
						Children: []any{
							textField("chapter_title", "Chapter title *", "chapter_title", "Chapter title", chapterTitle, "", true),
							textareaField("content", "Chapter *", "content", "Write your chapter...", chapterText(chapter), "Use {https://example.com/music.mp3} and {stop} for music cues.", 16, true),
							div("form-actions",
								elements.Button(constellation.WebComponent{
									Class:      "btn btn-primary",
									Attributes: map[string]string{"type": "submit"},
									Text:       "Save changes",
								}),
								buttonLink(chapterURL, "btn btn-outline", "Cancel"),
							),
						},
					}),
				),
				previewPanel("content"),
			),
			div("mt-10", buttonLink("/compose/chapter?novel="+url.QueryEscape(novel.ID), "text-link", "Add another chapter →")),
		),
		previewScript(),
	), http.StatusOK
}

func libraryPage(r *constellation.Request) string {
	novels, err := ListNovels()
	cards := make([]any, 0)

	if err == nil {
		for _, novel := range novels {
			description := strings.TrimSpace(novel.Description)
			if description == "" {
				description = "A story waiting to be read."
			}
			author := strings.TrimSpace(novel.Author)
			if author == "" {
				author = "Unknown author"
			}
			chapterWord := "chapters"
			if len(novel.Chapters) == 1 {
				chapterWord = "chapter"
			}

			cards = append(cards, elements.A(constellation.WebComponent{
				Class:      "card library-card",
				Attributes: map[string]string{"href": "/read/" + url.PathEscape(novel.ID)},
				Children: []any{
					span("card-eyebrow", strconv.Itoa(len(novel.Chapters))+" "+chapterWord),
					heading("h2", "card-title", novel.Title),
					paragraph("card-text", description),
					div("library-byline", span("", "by "+author), span("", "Open book →")),
				},
			}))
		}
	}

	if len(cards) == 0 {
		cards = append(cards, panel("empty-state p-10 text-center",
			heading("h2", "card-title", "No stories yet."),
			paragraph("card-text", "The library is waiting for its first book."),
			div("mt-6", buttonLink("/compose", "btn-primary", "Write the first one")),
		))
	}

	return page(r,
		section("py-14",
			eyebrow("The library"),
			heading("h1", "display mt-5 text-5xl", "Stories, ready to open."),
			lead("Every manuscript you save lives here. Choose a cover and settle in for a chapter."),
			divider(),
			div("grid gap-5 md:grid-cols-2", cards...),
			div("mt-10", buttonLink("/compose", "text-link", "Compose a novel →")),
		),
	)
}

func readerPage(r *constellation.Request) (string, int) {
	id := pathValue(r, "id")
	novel, err := GetNovel(id)
	if err != nil {
		return notFoundPage(r, "This story could not be found.", "/read", "Back to the library"), http.StatusNotFound
	}
	if len(novel.Chapters) == 0 {
		return notFoundPage(r, "This story has no chapters yet.", "/read", "Back to the library"), http.StatusNotFound
	}

	number := 1
	if raw := strings.TrimSpace(r.URL.Query().Get("chapter")); raw != "" {
		if parsed, parseErr := strconv.Atoi(raw); parseErr == nil && parsed > 0 {
			number = parsed
		}
	}
	if number > len(novel.Chapters) {
		number = len(novel.Chapters)
	}

	chapter := novel.Chapters[number-1]
	storyPath := "/read/" + url.PathEscape(novel.ID)
	chapterTitle := strings.TrimSpace(chapter.Title)
	if chapterTitle == "" {
		chapterTitle = "Chapter " + strconv.Itoa(number)
	}
	author := strings.TrimSpace(novel.Author)
	if author == "" {
		author = "Unknown author"
	}
	description := strings.TrimSpace(novel.Description)
	if description == "" {
		description = "A story from the Clio library."
	}

	paragraphs := make([]any, 0, len(chapter.Paragraphs))
	hasAudio := false
	for _, item := range chapter.Paragraphs {
		rendered := RenderMarkdown(item.Content)
		if strings.Contains(rendered, "data-audio-src") {
			hasAudio = true
		}
		paragraphs = append(paragraphs, div("reader-body", rendered))
	}

	tocItems := make([]any, 0, len(novel.Chapters))
	for index, item := range novel.Chapters {
		itemNumber := index + 1
		class := "toc-link"
		if itemNumber == number {
			class += " active"
		}
		title := strings.TrimSpace(item.Title)
		if title == "" {
			title = "Chapter " + strconv.Itoa(itemNumber)
		}
		tocItems = append(tocItems, elements.Li(constellation.WebComponent{
			Children: []any{link(storyPath+"?chapter="+strconv.Itoa(itemNumber), class, title)},
		}))
	}

	navigation := make([]any, 0, 2)
	if number > 1 {
		navigation = append(navigation, buttonLink(storyPath+"?chapter="+strconv.Itoa(number-1), "btn-outline", "← Previous chapter"))
	}
	if number < len(novel.Chapters) {
		navigation = append(navigation, buttonLink(storyPath+"?chapter="+strconv.Itoa(number+1), "btn-primary", "Next chapter →"))
	}

	children := make([]any, 0, 3)
	if hasAudio {
		children = append(children,
			elements.Div(constellation.WebComponent{
				Class:      "start-overlay",
				Attributes: map[string]string{"id": "start-overlay"},
				Children: []any{
					div("start-card",
						eyebrow("Before you begin"),
						heading("h2", "display mt-4 text-3xl", novel.Title),
						paragraph("mt-4 text-sm", "Start reading to allow background music to play as you scroll."),
						elements.Button(constellation.WebComponent{
							Class:      "btn btn-primary mt-7",
							Attributes: map[string]string{"id": "start-reading", "type": "button"},
							Text:       "Start reading",
						}),
					),
				},
			}),
			elements.Audio(constellation.WebComponent{
				Attributes: map[string]string{"id": "story-audio", "loop": "loop", "preload": "none"},
			}),
			elements.Div(constellation.WebComponent{
				Class:      "now-playing",
				Attributes: map[string]string{"id": "now-playing", "hidden": "hidden"},
				Children: []any{
					span("now-playing-pulse", "♪"),
					span("now-playing-label", "Now playing"),
					span("now-title", ""),
					elements.Button(constellation.WebComponent{
						Class:      "player-btn",
						Attributes: map[string]string{"id": "audio-toggle", "type": "button", "aria-label": "Pause music"},
						Text:       "❚❚",
					}),
					elements.Button(constellation.WebComponent{
						Class:      "player-btn",
						Attributes: map[string]string{"id": "audio-mute", "type": "button", "aria-label": "Mute music"},
						Text:       "Mute",
					}),
				},
			}),
		)
	}

	children = append(children,
		section("py-14",
			buttonLink("/read", "back-link", "← Back to the library"),
			div("mt-10 space-y-4",
				eyebrow("Now reading"),
				heading("h1", "display text-4xl sm:text-5xl", novel.Title),
				paragraph("byline", "by "+author),
				paragraph("lead", description),
			),
			divider(),
			div("grid items-start gap-12 lg:grid-cols-[1.35fr_0.65fr]",
				div("",
					div("reader-chapter-head px-5",
						span("eyebrow", "Chapter "+strconv.Itoa(number)+" of "+strconv.Itoa(len(novel.Chapters))),
						heading("h2", "display mt-3 text-3xl", chapterTitle),
					),
					div("reader-stack mt-8", paragraphs...),
					div("reader-nav mt-10", navigation...),
				),
				panel("side-panel p-6",
					span("eyebrow", "Contents"),
					elements.Ol(constellation.WebComponent{Class: "toc mt-5", Children: tocItems}),
					div("side-actions mt-7",
						buttonLink("/compose/chapter?novel="+url.PathEscape(novel.ID), "btn btn-outline w-full", "Add a chapter"),
						buttonLink("/compose/chapter/"+url.PathEscape(novel.ID)+"/"+strconv.Itoa(number)+"/edit", "btn btn-outline w-full", "Edit this chapter"),
					),
				),
			),
		),
	)
	if hasAudio {
		children = append(children, readerScript())
	}

	return page(r, children...), http.StatusOK
}
