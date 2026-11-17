package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const novelsFile = "src/backend/novels.json"

type Paragraph struct {
	Content string `json:"content"`
}

type Chapter struct {
	Title      string      `json:"title"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Novel struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Chapters    []Chapter `json:"chapters"`
}

type novelDatabase struct {
	Novels []Novel `json:"novels"`
}

type NovelInput struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Description  string `json:"description"`
	ChapterTitle string `json:"chapter_title"`
	Content      string `json:"content"`
}

type ChapterInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	storeMu      sync.Mutex
	paragraphSep = regexp.MustCompile(`\n\s*\n`)
	nonAlphaNum  = regexp.MustCompile(`[^a-z0-9]+`)
)

var (
	ErrNovelNotFound   = errors.New("novel not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrInvalidInput    = errors.New("title and content are required")
)

func readDatabase() (novelDatabase, error) {
	data, err := os.ReadFile(novelsFile)
	if err != nil {
		return novelDatabase{}, fmt.Errorf("read novels: %w", err)
	}

	var db novelDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		return novelDatabase{}, fmt.Errorf("parse novels: %w", err)
	}
	return db, nil
}

func writeDatabase(db novelDatabase) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("encode novels: %w", err)
	}
	if err := os.WriteFile(novelsFile, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write novels: %w", err)
	}
	return nil
}

func splitParagraphs(content string) []Paragraph {
	parts := paragraphSep.Split(strings.TrimSpace(content), -1)
	paragraphs := make([]Paragraph, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			paragraphs = append(paragraphs, Paragraph{Content: part})
		}
	}
	return paragraphs
}

func slug(title string) string {
	s := nonAlphaNum.ReplaceAllString(strings.ToLower(strings.TrimSpace(title)), "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "story"
	}
	return fmt.Sprintf("%s-%d", s, time.Now().UnixMilli())
}

// ListNovels returns every saved novel.
func ListNovels() ([]Novel, error) {
	storeMu.Lock()
	defer storeMu.Unlock()

	db, err := readDatabase()
	if err != nil {
		return nil, err
	}
	return db.Novels, nil
}

// GetNovel finds a novel by its id.
func GetNovel(id string) (Novel, error) {
	storeMu.Lock()
	defer storeMu.Unlock()

	db, err := readDatabase()
	if err != nil {
		return Novel{}, err
	}
	for _, novel := range db.Novels {
		if novel.ID == id {
			return novel, nil
		}
	}
	return Novel{}, ErrNovelNotFound
}

// CreateNovel saves a new novel with its first chapter and returns it.
func CreateNovel(in NovelInput) (Novel, error) {
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" || content == "" {
		return Novel{}, ErrInvalidInput
	}

	chapterTitle := strings.TrimSpace(in.ChapterTitle)
	if chapterTitle == "" {
		chapterTitle = "Chapter 1"
	}

	novel := Novel{
		ID:          slug(title),
		Title:       title,
		Author:      strings.TrimSpace(in.Author),
		Description: strings.TrimSpace(in.Description),
		Chapters: []Chapter{{
			Title:      chapterTitle,
			Paragraphs: splitParagraphs(content),
		}},
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	db, err := readDatabase()
	if err != nil {
		return Novel{}, err
	}
	db.Novels = append(db.Novels, novel)
	if err := writeDatabase(db); err != nil {
		return Novel{}, err
	}
	return novel, nil
}

// AddChapter appends a chapter to an existing novel and returns the novel.
func AddChapter(id string, in ChapterInput) (Novel, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return Novel{}, ErrInvalidInput
	}

	chapterTitle := strings.TrimSpace(in.Title)
	autoTitle := chapterTitle == ""

	storeMu.Lock()
	defer storeMu.Unlock()

	db, err := readDatabase()
	if err != nil {
		return Novel{}, err
	}

	for i := range db.Novels {
		if db.Novels[i].ID == id {
			if autoTitle {
				chapterTitle = fmt.Sprintf("Chapter %d", len(db.Novels[i].Chapters)+1)
			}
			db.Novels[i].Chapters = append(db.Novels[i].Chapters, Chapter{
				Title:      chapterTitle,
				Paragraphs: splitParagraphs(content),
			})
			if err := writeDatabase(db); err != nil {
				return Novel{}, err
			}
			return db.Novels[i], nil
		}
	}
	return Novel{}, ErrNovelNotFound
}

// UpdateChapter replaces the title and text of a chapter (1-based chapter number).
func UpdateChapter(id string, chapterNumber int, in ChapterInput) (Novel, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return Novel{}, ErrInvalidInput
	}

	chapterTitle := strings.TrimSpace(in.Title)
	if chapterTitle == "" {
		chapterTitle = fmt.Sprintf("Chapter %d", chapterNumber)
	}

	storeMu.Lock()
	defer storeMu.Unlock()

	db, err := readDatabase()
	if err != nil {
		return Novel{}, err
	}

	for i := range db.Novels {
		if db.Novels[i].ID != id {
			continue
		}
		if chapterNumber < 1 || chapterNumber > len(db.Novels[i].Chapters) {
			return Novel{}, ErrChapterNotFound
		}
		db.Novels[i].Chapters[chapterNumber-1] = Chapter{
			Title:      chapterTitle,
			Paragraphs: splitParagraphs(content),
		}
		if err := writeDatabase(db); err != nil {
			return Novel{}, err
		}
		return db.Novels[i], nil
	}
	return Novel{}, ErrNovelNotFound
}
