# Clio

Clio is a small web-based writing studio for designing novels. It lets you draft
chapters in Markdown, preview them as you type, save them in a local library, and
read them back with book-like typography and optional background-music markers.

It is a Go application built with [Constellation](https://github.com/RyannKim327/constellation),
a convention-over-configuration web framework built on Gin. Pages are
server-rendered, static assets are embedded into the binary, and story data is
stored in a JSON file.

> Clio was used as the test subject for Constellation, the framework created by
> Ryann Kim Sesgundo.

## Features

- **Novel composer** — create a novel with a title, author, description, and
  first chapter from `/compose`.
- **Chapter management** — append chapters and edit existing ones; blank chapter
  titles are numbered automatically.
- **Live Markdown preview** — writing pages post to `/api/markdown` and render a
  preview with GitHub Flavored Markdown via Goldmark.
- **Reader** — browse the library at `/read`, then read one chapter at a time with
  a table of contents, previous/next navigation, and `?chapter=N` deep links.
- **Story audio markers** — write `{https://example.com/music.mp3}` in a chapter
  to start background music and `{stop}` to stop it. The reader unlocks playback
  with a "Start reading" overlay, shows a now-playing bar, and pauses when the
  tab is hidden.
- **JSON API** — list and fetch novels, create novels, add/update chapters, and
  render Markdown. Form submissions and JSON requests are both supported.
- **Batteries included** — CORS, basic bot/scraper detection, optional cookie
  access guard, request logging, custom 404 page, and embedded static files.

## Requirements

- Go **1.27+**
- (optional) [Air](https://github.com/air-verse/air) for live reload

## Getting started

Run commands from the repository root; the JSON store is resolved relative to it.

```bash
git clone <repository-url>
cd istorya
go mod download
go run .
```

Then open <http://localhost:8000>.

You can also run the binary directly:

```bash
go build -o clio ./main.go
./clio
```

### Development with live reload

```bash
air
```

Air uses `.air.toml`, builds `main.go` into `tmp/main`, and rebuilds when `.go`,
`.css`, `.js`, or `.html` files change.

## Configuration

| Variable          | Default                                                                 | Description                                  |
| ----------------- | ----------------------------------------------------------------------- | -------------------------------------------- |
| `PORT`            | `8000`                                                                  | Port the HTTP server listens on.             |
| `ALLOWED_ORIGINS` | `https://constellation-red.vercel.app,http://localhost:3000,http://localhost:5173` | Comma-separated CORS allow-list. |

`PORT` takes precedence over the programmatic port argument, which falls back to
`8000`. CORS allows `GET`, `POST`, `PATCH`, `PUT`, `DELETE`, and `OPTIONS`.

## Routes

### Pages

| Method | Path                                  | Description                                      |
| ------ | ------------------------------------- | ------------------------------------------------ |
| `GET`  | `/`                                   | Landing page.                                    |
| `GET`  | `/about`                              | Project overview.                                |
| `GET`  | `/compose`                            | New-novel form with live Markdown preview.       |
| `GET`  | `/compose/chapter`                    | Add a chapter to an existing novel.              |
| `GET`  | `/compose/chapter/:id/:chapter/edit`  | Edit an existing chapter.                        |
| `GET`  | `/read`                               | Novel library.                                   |
| `GET`  | `/read/:id`                           | Chapter reader with TOC and audio controls.      |
| `GET`  | `/static/*`                           | Embedded stylesheet and other static assets.     |

### API

| Method | Path                                        | Description                        |
| ------ | ------------------------------------------- | ---------------------------------- |
| `GET`  | `/api`                                      | Example JSON response.             |
| `GET`  | `/api/novels`                               | List all novels.                   |
| `GET`  | `/api/novels/:id`                           | Get one novel.                     |
| `POST` | `/api/novels`                               | Create a novel with its first chapter. |
| `POST` | `/api/novels/:id/chapters`                  | Append a chapter.                  |
| `PUT`  | `/api/novels/:id/chapters/:chapter`         | Update a chapter (1-based number). |
| `POST` | `/api/novels/:id/chapters/:chapter/update`  | Form-friendly update endpoint.     |
| `POST` | `/api/markdown`                             | Render story Markdown to HTML.     |

API handlers return JSON when the request `Content-Type` is `application/json`;
otherwise they accept form fields and redirect the browser after a successful
write.

## Data storage

Novels are stored in `src/backend/novels.json` as a single `novels` array. Each
novel has an ID, title, author, description, and chapters; each chapter has a
title and paragraphs of Markdown content. Access is serialized with a mutex, and
IDs are generated from a slug plus a millisecond timestamp.

This is intentionally simple file-backed storage: there is no database,
authentication system, or user accounts yet.

## Architecture

`main.go` calls `src.Main`, which registers the routes in `src/index.go` and
hands them to the bootstrap layer in `config/main.go`.

A route is described with `middleware.Route`:

```go
var Page = middleware.Route{
    Method:     utils.METHOD_GET,
    Path:       "/example",
    Permission: utils.PERMISSION_ALL,
    Title:      "Example · Clio",
    Handler:    handler,
}
```

Handlers return an HTML body plus a parameter map, and the framework wraps the
result in the page shell. API handlers write directly with `ctx.JSON(...)` and
return an empty body with `nil` parameters.

## Project structure

```
.
├── main.go                    # Entry point
├── config/
│   ├── main.go                # Bootstrap: CORS, route registration, server start
│   ├── middleware/            # Routes, CORS, bot guard, cookie guard, server
│   └── utils/                 # Templating, rendering, constants
├── src/
│   ├── index.go               # App route registry
│   ├── embed.go               # Embedded static assets
│   ├── backend/               # Novel store and JSON APIs
│   │   └── novels.json        # File-backed novel data
│   ├── routes/                # Page handlers
│   ├── components/            # Shared header and live-preview components
│   ├── utils/                 # Story Markdown rendering
│   └── static/styles.css      # Site stylesheet
├── tools/md2go/               # Standalone Markdown-to-Go generator
└── .air.toml                  # Air live-reload configuration
```

## Notes

- Bot protection is a simple `User-Agent` heuristic; harden it before relying on
  it in production.
- The cookie access guard exists in the framework but is not used by Clio's
  current routes.
- Some lower-level package and module names still use earlier project names; the
  product name is Clio.

## License

No license has been specified for this project.

## Credits

> [!NOTE]
> People mentioned here are not just part of the project, but were also part of my journey.

**Ground Zero Community**
* Justeen James Tolentino
* Salvador
* John Roy Lapida
* Jay Patrick Cano
* John Paul Caigas

**Droid Modificators**
* John Jeremy Antigou
* Mark Kevin Manalo
* Adrian Bundas
* Genesis Fernandez

**Philtech Developers**
* Fortzkie Tuneskie
* Dominic Ama
