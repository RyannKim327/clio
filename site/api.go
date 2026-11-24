package site

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	constellation "constellation/src"
	"constellation/src/middleware"
)

type novelRequest struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Description  string `json:"description"`
	ChapterTitle string `json:"chapter_title"`
	Content      string `json:"content"`
}

type chapterRequest struct {
	Title        string `json:"title"`
	ChapterTitle string `json:"chapter_title"`
	Content      string `json:"content"`
}

func httpRequest(r *constellation.Request) *http.Request {
	return (*http.Request)(r)
}

func pathValue(r *constellation.Request, name string) string {
	return httpRequest(r).PathValue(name)
}

func wantsJSON(r *constellation.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func decodeNovelInput(r *constellation.Request) (NovelInput, error) {
	var payload novelRequest

	if wantsJSON(r) {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return NovelInput{}, err
		}
	} else {
		if err := httpRequest(r).ParseForm(); err != nil {
			return NovelInput{}, err
		}
		payload.Title = httpRequest(r).FormValue("title")
		payload.Author = httpRequest(r).FormValue("author")
		payload.Description = httpRequest(r).FormValue("description")
		payload.ChapterTitle = httpRequest(r).FormValue("chapter_title")
		payload.Content = httpRequest(r).FormValue("content")
	}

	return NovelInput{
		Title:        payload.Title,
		Author:       payload.Author,
		Description:  payload.Description,
		ChapterTitle: payload.ChapterTitle,
		Content:      payload.Content,
	}, nil
}

func decodeChapterInput(r *constellation.Request) (ChapterInput, error) {
	var payload chapterRequest

	if wantsJSON(r) {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return ChapterInput{}, err
		}
	} else {
		if err := httpRequest(r).ParseForm(); err != nil {
			return ChapterInput{}, err
		}
		payload.Title = httpRequest(r).FormValue("chapter_title")
		payload.Content = httpRequest(r).FormValue("content")
	}

	title := strings.TrimSpace(payload.Title)
	if title == "" {
		title = strings.TrimSpace(payload.ChapterTitle)
	}
	return ChapterInput{Title: title, Content: payload.Content}, nil
}

func writeAPI(w constellation.Response, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	middleware.SendApi(w, payload)
}

func writeAPIError(w constellation.Response, status int, message string) {
	if status >= http.StatusInternalServerError {
		message = "internal server error"
	}
	writeAPI(w, status, map[string]any{"error": message})
}

func limitBody(r *constellation.Request, w constellation.Response) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
}

func writeDecodeError(w constellation.Response, err error) {
	if strings.Contains(err.Error(), "request body too large") {
		writeAPIError(w, http.StatusRequestEntityTooLarge, "request body too large")
		return
	}
	writeDecodeError(w, err)
}

func writeStoreError(w constellation.Response, err error) {
	status := statusFor(err)
	if status >= http.StatusInternalServerError {
		log.Printf("store error: %v", err)
	}
	writeAPIError(w, status, err.Error())
}

func chapterFromPath(r *constellation.Request) (int, bool) {
	number, err := strconv.Atoi(pathValue(r, "chapter"))
	return number, err == nil && number > 0
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, ErrNovelNotFound), errors.Is(err, ErrChapterNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func apiIndex(r *constellation.Request, w constellation.Response) {
	writeAPI(w, http.StatusOK, map[string]any{"message": "Hello"})
}

func apiNotFound(r *constellation.Request, w constellation.Response) {
	http.NotFound(w, httpRequest(r))
}

func writePage(w constellation.Response, title, html string, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	middleware.WebParser(w, title, html)
}

func readerEndpoint(r *constellation.Request, w constellation.Response) {
	html, status := readerPage(r)
	writePage(w, "Read · Clio", html, status)
}

func editChapterEndpoint(r *constellation.Request, w constellation.Response) {
	html, status := editChapterPage(r)
	writePage(w, "Edit chapter · Clio", html, status)
}

func apiListNovels(r *constellation.Request, w constellation.Response) {
	novels, err := ListNovels()
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeAPI(w, http.StatusOK, map[string]any{"novels": novels})
}

func apiGetNovel(r *constellation.Request, w constellation.Response) {
	novel, err := GetNovel(pathValue(r, "id"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeAPI(w, http.StatusOK, map[string]any{"novel": novel})
}

func apiCreateNovel(r *constellation.Request, w constellation.Response) {
	limitBody(r, w)
	input, err := decodeNovelInput(r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	novel, err := CreateNovel(input)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if wantsJSON(r) {
		writeAPI(w, http.StatusCreated, map[string]any{"novel": novel})
		return
	}
	http.Redirect(w, httpRequest(r), "/read/"+novel.ID, http.StatusSeeOther)
}

func apiAddChapter(r *constellation.Request, w constellation.Response) {
	limitBody(r, w)
	input, err := decodeChapterInput(r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	id := pathValue(r, "id")
	if !wantsJSON(r) {
		if selected := strings.TrimSpace(httpRequest(r).FormValue("novel_id")); selected != "" {
			id = selected
		}
	}
	novel, err := AddChapter(id, input)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if wantsJSON(r) {
		writeAPI(w, http.StatusCreated, map[string]any{"novel": novel})
		return
	}
	http.Redirect(w, httpRequest(r), "/read/"+novel.ID+"?chapter="+strconv.Itoa(len(novel.Chapters)), http.StatusSeeOther)
}

func apiUpdateChapter(r *constellation.Request, w constellation.Response) {
	limitBody(r, w)
	number, valid := chapterFromPath(r)
	if !valid {
		writeAPIError(w, http.StatusBadRequest, "invalid chapter number")
		return
	}

	input, err := decodeChapterInput(r)
	if err != nil {
		writeDecodeError(w, err)
		return
	}

	novel, err := UpdateChapter(pathValue(r, "id"), number, input)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if wantsJSON(r) {
		writeAPI(w, http.StatusOK, map[string]any{"novel": novel})
		return
	}
	http.Redirect(w, httpRequest(r), "/read/"+novel.ID+"?chapter="+strconv.Itoa(number), http.StatusSeeOther)
}

func apiMarkdown(r *constellation.Request, w constellation.Response) {
	limitBody(r, w)
	var text string

	if wantsJSON(r) {
		var payload struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeDecodeError(w, err)
			return
		}
		text = payload.Text
	} else {
		if err := httpRequest(r).ParseForm(); err != nil {
			writeDecodeError(w, err)
			return
		}
		text = httpRequest(r).FormValue("text")
	}

	writeAPI(w, http.StatusOK, map[string]any{"html": RenderMarkdown(text)})
}

func staticFile(assets fs.FS) constellation.HandlerApi {
	return func(r *constellation.Request, w constellation.Response) {
		name := pathValue(r, "file")
		if name == "" || path.Base(name) != name || strings.Contains(name, "..") {
			http.NotFound(w, httpRequest(r))
			return
		}

		data, err := fs.ReadFile(assets, "static/"+name)
		if err != nil {
			http.NotFound(w, httpRequest(r))
			return
		}

		contentType := mime.TypeByExtension(path.Ext(name))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(data)
	}
}
