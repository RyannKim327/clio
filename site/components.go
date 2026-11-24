package site

import (
	"strconv"
	"strings"

	constellation "constellation/src"
	"constellation/src/elements"
)

func page(r *constellation.Request, content ...any) string {
	children := []any{siteHeader(activeSection(r.URL.Path))}
	children = append(children, content...)
	children = append(children, siteFooter())
	return elements.Div(constellation.WebComponent{
		Class:    "site-shell mx-auto w-full max-w-6xl px-5 pb-16 sm:px-8",
		Children: children,
	})
}

func activeSection(path string) string {
	value := strings.TrimPrefix(path, "/")
	switch {
	case value == "":
		return "home"
	case value == "compose" || strings.HasPrefix(value, "compose/"):
		return "compose"
	case value == "read" || strings.HasPrefix(value, "read/"):
		return "read"
	case value == "about" || strings.HasPrefix(value, "about/"):
		return "about"
	default:
		return ""
	}
}

func siteHeader(active string) string {
	items := []struct {
		name string
		href string
		key  string
	}{
		{name: "Home", href: "/", key: "home"},
		{name: "Compose", href: "/compose", key: "compose"},
		{name: "Read", href: "/read", key: "read"},
		{name: "About", href: "/about", key: "about"},
	}

	links := make([]any, 0, len(items))
	for _, item := range items {
		class := "nav-link"
		if item.key == active {
			class += " active"
		}
		links = append(links, link(item.href, class, item.name))
	}

	return elements.Header(constellation.WebComponent{
		Class: "site-header flex flex-wrap items-center justify-between gap-5 border-b border-white/10 py-6 px-10",
		Children: []any{
			elements.A(constellation.WebComponent{
				Class:      "brand flex items-center gap-3",
				Attributes: map[string]string{"href": "/", "aria-label": "Clio home"},
				Children: []any{
					clioLogo("h-8 w-8 shrink-0"),
					elements.Span(constellation.WebComponent{Text: "Clio"}),
					elements.Span(constellation.WebComponent{Class: "brand-dot", Text: "."}),
				},
			}),
			elements.Nav(constellation.WebComponent{
				Class:    "site-nav flex flex-wrap items-center gap-6",
				Children: links,
			}),
		},
	})
}

func siteFooter() string {
	return elements.Footer(constellation.WebComponent{
		Class: "site-footer mt-20 flex flex-wrap items-center justify-between gap-4 border-t border-white/10 pt-7 text-sm text-slate-500",
		Children: []any{
			elements.Span(constellation.WebComponent{Text: "Clio · a home for your stories"}),
			link("/about", "footer-link", "About the project"),
		},
	})
}

func link(href, class, text string) string {
	return elements.A(constellation.WebComponent{
		Class:      class,
		Attributes: map[string]string{"href": href},
		Text:       text,
	})
}

func buttonLink(href, class, text string) string {
	return link(href, "btn "+class, text)
}

func div(class string, children ...any) string {
	return elements.Div(constellation.WebComponent{Class: class, Children: children})
}

func section(class string, children ...any) string {
	return elements.Section(constellation.WebComponent{Class: class, Children: children})
}

func span(class, text string) string {
	return elements.Span(constellation.WebComponent{Class: class, Text: text})
}

func paragraph(class, text string) string {
	return elements.P(constellation.WebComponent{Class: class, Text: text})
}

func heading(tag, class, text string) string {
	component := constellation.WebComponent{Class: class, Text: text}
	switch tag {
	case "h1":
		return elements.H1(component)
	case "h2":
		return elements.H2(component)
	case "h3":
		return elements.H3(component)
	default:
		return elements.H4(component)
	}
}

func eyebrow(text string) string {
	return span("eyebrow", text)
}

func lead(text string) string {
	return paragraph("lead max-w-2xl", text)
}

func divider() string {
	return div("divider my-14", span("divider-mark", "✦"))
}

func panel(class string, children ...any) string {
	return div("panel "+class, children...)
}

func textField(id, label, name, placeholder, value, hint string, required bool) string {
	attributes := map[string]string{
		"id":          id,
		"name":        name,
		"type":        "text",
		"placeholder": placeholder,
		"value":       value,
	}
	if required {
		attributes["required"] = "required"
	}

	children := []any{
		elements.Label(constellation.WebComponent{
			Class:      "field-label",
			Attributes: map[string]string{"for": id},
			Text:       label,
		}),
		elements.Input(constellation.WebComponent{Class: "field-input", Attributes: attributes}),
	}
	if hint != "" {
		children = append(children, paragraph("form-hint", hint))
	}
	return div("field", children...)
}

func textareaField(id, label, name, placeholder, value, hint string, rows int, required bool) string {
	attributes := map[string]string{
		"id":          id,
		"name":        name,
		"placeholder": placeholder,
		"rows":        strconv.Itoa(rows),
	}
	if required {
		attributes["required"] = "required"
	}

	children := []any{
		elements.Label(constellation.WebComponent{
			Class:      "field-label",
			Attributes: map[string]string{"for": id},
			Text:       label,
		}),
		elements.Textarea(constellation.WebComponent{Class: "field-input field-textarea", Attributes: attributes, Text: value}),
	}
	if hint != "" {
		children = append(children, paragraph("form-hint", hint))
	}
	return div("field", children...)
}

func selectField(id, label, name string, options ...any) string {
	return div("field",
		elements.Label(constellation.WebComponent{
			Class:      "field-label",
			Attributes: map[string]string{"for": id},
			Text:       label,
		}),
		elements.Select(constellation.WebComponent{
			Class: "field-input",
			Attributes: map[string]string{
				"id":       id,
				"name":     name,
				"required": "required",
			},
			Children: options,
		}),
	)
}

func previewPanel(inputID string) string {
	return div("preview-panel",
		div("preview-heading", span("eyebrow", "Live preview"), span("preview-live", "● live")),
		elements.Div(constellation.WebComponent{
			Class: "preview-body story-body",
			Attributes: map[string]string{
				"data-preview-for": inputID,
				"aria-live":        "polite",
			},
		}),
	)
}

func previewScript() string {
	return `<script>
(function () {
  var timers = {};
  document.querySelectorAll('[data-preview-for]').forEach(function (panel) {
    var input = document.getElementById(panel.getAttribute('data-preview-for'));
    if (!input) return;
    function render() {
      fetch('/api/markdown', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: input.value })
      }).then(function (response) {
        return response.json();
      }).then(function (data) {
        panel.innerHTML = data.html || '';
      }).catch(function () {
        panel.textContent = 'Preview unavailable.';
      });
    }
    input.addEventListener('input', function () {
      clearTimeout(timers[input.id]);
      timers[input.id] = setTimeout(render, 150);
    });
    render();
  });
})();
</script>`
}

func chapterSwitchScript() string {
	return `<script>
(function () {
  var select = document.getElementById('novel-select');
  var form = document.getElementById('chapter-form');
  if (!select || !form) return;
  select.addEventListener('change', function () {
    form.action = '/api/novels/' + encodeURIComponent(this.value) + '/chapters';
  });
})();
</script>`
}

func readerScript() string {
	return `<script>
(function () {
  var overlay = document.getElementById('start-overlay');
  var start = document.getElementById('start-reading');
  var audio = document.getElementById('story-audio');
  var bar = document.getElementById('now-playing');
  var title = document.getElementById('now-title');
  var toggle = document.getElementById('audio-toggle');
  var mute = document.getElementById('audio-mute');
  if (!overlay || !start || !audio || !bar || !title || !toggle || !mute) return;

  var unlocked = false;
  var current = '';

  function trackName(address) {
    try {
      var parsed = new URL(address, window.location.href);
      var segment = parsed.pathname.split('/').filter(Boolean).pop() || parsed.hostname;
      return decodeURIComponent(segment);
    } catch (error) {
      return 'Background music';
    }
  }

  function play() {
    var result = audio.play();
    if (result && result.catch) result.catch(function () {});
  }

  function setTrack(address) {
    if (!address || address === current) return;
    current = address;
    audio.src = address;
    audio.loop = true;
    title.textContent = trackName(address);
    bar.hidden = false;
    play();
  }

  function stop() {
    audio.pause();
    audio.currentTime = 0;
    current = '';
    bar.hidden = true;
  }

  start.addEventListener('click', function () {
    unlocked = true;
    overlay.hidden = true;
    var first = document.querySelector('[data-audio-src]');
    if (first) setTrack(first.getAttribute('data-audio-src'));
  });

  toggle.addEventListener('click', function () {
    if (audio.paused) play(); else audio.pause();
  });

  mute.addEventListener('click', function () {
    audio.muted = !audio.muted;
    mute.textContent = audio.muted ? 'Unmute' : 'Mute';
  });

  document.addEventListener('visibilitychange', function () {
    if (document.hidden) audio.pause();
    else if (unlocked && audio.src) play();
  });

  if ('IntersectionObserver' in window) {
    var observer = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (!entry.isIntersecting || !unlocked) return;
        if (entry.target.hasAttribute('data-audio-stop')) stop();
        else setTrack(entry.target.getAttribute('data-audio-src'));
      });
    }, { threshold: 0.5 });
    document.querySelectorAll('[data-audio-src],[data-audio-stop]').forEach(function (marker) {
      observer.observe(marker);
    });
  }
})();
</script>`
}

func notFoundPage(r *constellation.Request, message, backHref, backText string) string {
	return page(r,
		section("py-24 text-center",
			eyebrow("Not found"),
			heading("h1", "display mt-5 text-4xl", message),
			div("mt-8", buttonLink(backHref, "btn-primary", backText)),
		),
	)
}
