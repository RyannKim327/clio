# Clio

Clio is a small web-based writing studio for designing novels. Draft chapters in
Markdown, preview them as you type, keep everything in a local JSON library, and
read chapters back with book-like typography and optional music markers.

The site runs on the fixed Constellation framework in `src/`. All application code
lives outside `src/` in `site/`, with `main.go` wiring the endpoints together.

> The `src/` directory must never be modified by agents. It may be read, but all
> changes belong in the rest of the repository. See [AGENTS.md](AGENTS.md).

## Features

- **Novel composer** — create a novel with a title, author, description, and
  first chapter at `/compose`.
- **Chapter management** — append chapters at `/compose/chapter` and edit them
  from the reader. Blank chapter titles are numbered automatically.
- **Live Markdown preview** — writing pages post to `/api/markdown` and render a
  preview built with the standard library only.
- **Reader** — browse `/read`, then read one chapter at a time with a table of
  contents, previous/next navigation, and `?chapter=N` deep links.
- **Story audio markers** — write `{https://example.com/music.mp3}` in a chapter
  to start background music and `{stop}` to stop it. The reader unlocks playback
  with a "Start reading" overlay and pauses when the tab is hidden.
- **JSON API** — list and fetch novels, create novels, add/update chapters, and
  render Markdown. Form submissions and JSON requests are both supported.

## Requirements

- Go 1.27 or newer

## Getting started

Run from the repository root so the default data path resolves correctly.

```bash
go run .
```

Open <http://localhost:8000>.

Build a binary with `go build -o clio ./main.go`.

### Development with live reload

```bash
air
```

`.air.toml` builds `main.go` into `tmp/main` and rebuilds when relevant files
change.

## Configuration

| Variable         | Default             | Description                              |
| ---------------- | ------------------- | ---------------------------------------- |
| `PORT`           | `8000`              | Port the HTTP server listens on.         |
| `CLIO_DATA_FILE` | `data/novels.json`  | Path to the JSON story library.          |

## Routes

### Pages

| Method | Path                                    | Description                                |
| ------ | --------------------------------------- | ------------------------------------------ |
| `GET`  | `/`                                     | Landing page.                              |
| `GET`  | `/about`                                | Project overview.                          |
| `GET`  | `/compose`                              | New-novel form with live Markdown preview. |
| `GET`  | `/compose/chapter`                      | Add a chapter to an existing novel.        |
| `GET`  | `/compose/chapter/{id}/{chapter}/edit`  | Edit an existing chapter.                  |
| `GET`  | `/read`                                 | Novel library.                             |
| `GET`  | `/read/{id}`                            | Chapter reader with TOC and audio.         |
| `GET`  | `/static/{file}`                        | Embedded stylesheet and logo.              |

### API

| Method | Path                                        | Description                              |
| ------ | ------------------------------------------- | ---------------------------------------- |
| `GET`  | `/api`                                      | Example JSON response.                   |
| `GET`  | `/api/novels`                               | List all novels.                         |
| `GET`  | `/api/novels/{id}`                          | Get one novel.                           |
| `POST` | `/api/novels`                               | Create a novel with its first chapter.   |
| `POST` | `/api/novels/{id}/chapters`                 | Append a chapter.                        |
| `PUT`  | `/api/novels/{id}/chapters/{chapter}`       | Update a chapter (1-based number).       |
| `POST` | `/api/novels/{id}/chapters/{chapter}/update`| Form-friendly update endpoint.           |
| `POST` | `/api/markdown`                             | Render story Markdown to HTML.           |

JSON requests send `Content-Type: application/json`; otherwise handlers read
form fields and redirect with `303 See Other` after a successful write.

## Data storage

Novels are stored in `data/novels.json` as a `novels` array. Each novel has an
ID, title, author, description, and chapters; each chapter has a title and
paragraphs of Markdown source. Access is serialized with a mutex, and IDs are
generated from a slug plus a millisecond timestamp.

## Markdown

The renderer in `site/markdown.go` is intentionally dependency-free. It supports
headings, emphasis, strikethrough, highlights, quotes, lists, task items, code,
tables, links, images, and horizontal rules. Raw HTML is escaped and link schemes
are restricted to safe values.

## Verification

```bash
gofmt -l .
go vet ./...
go build ./...
go test ./...
```

## License

See [LICENSE.md](LICENSE.md).
