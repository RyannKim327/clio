package site

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const defaultDataFile = "data/novels.json"

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
	Title        string
	Author       string
	Description  string
	ChapterTitle string
	Content      string
}

type ChapterInput struct {
	Title   string
	Content string
}

var (
	ErrNovelNotFound   = errors.New("novel not found")
	ErrChapterNotFound = errors.New("chapter not found")
	ErrInvalidInput    = errors.New("title and content are required")

	storeMu      sync.RWMutex
	paragraphSep = regexp.MustCompile(`\n\s*\n`)
	nonAlphaNum  = regexp.MustCompile(`[^a-z0-9]+`)
)

func dataFile() string {
	if value := strings.TrimSpace(os.Getenv("CLIO_DATA_FILE")); value != "" {
		return value
	}
	return defaultDataFile
}

func readDatabase() (novelDatabase, error) {
	data, err := os.ReadFile(dataFile())
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

	path := dataFile()
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create data directory: %w", err)
		}
	}

	temp, err := os.CreateTemp(dir, ".novels-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)

	if _, err := temp.Write(append(data, '\n')); err != nil {
		temp.Close()
		return fmt.Errorf("write novels: %w", err)
	}
	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return fmt.Errorf("set novel permissions: %w", err)
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return fmt.Errorf("sync novels: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close novels: %w", err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("replace novels: %w", err)
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

func chapterText(chapter Chapter) string {
	parts := make([]string, 0, len(chapter.Paragraphs))
	for _, paragraph := range chapter.Paragraphs {
		parts = append(parts, strings.TrimSpace(paragraph.Content))
	}
	return strings.Join(parts, "\n\n")
}

func slug(title string) string {
	s := nonAlphaNum.ReplaceAllString(strings.ToLower(strings.TrimSpace(title)), "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "story"
	}
	return fmt.Sprintf("%s-%d", s, time.Now().UnixMilli())
}

func hasNovel(novels []Novel, id string) bool {
	for _, novel := range novels {
		if novel.ID == id {
			return true
		}
	}
	return false
}

func ListNovels() ([]Novel, error) {
	storeMu.RLock()
	defer storeMu.RUnlock()

	db, err := readDatabase()
	if err != nil {
		return nil, err
	}
	return db.Novels, nil
}

func GetNovel(id string) (Novel, error) {
	storeMu.RLock()
	defer storeMu.RUnlock()

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
	base := slug(title)
	novel.ID = base
	for suffix := 2; hasNovel(db.Novels, novel.ID); suffix++ {
		novel.ID = fmt.Sprintf("%s-%d", base, suffix)
	}
	db.Novels = append(db.Novels, novel)
	if err := writeDatabase(db); err != nil {
		return Novel{}, err
	}
	return novel, nil
}

func AddChapter(id string, in ChapterInput) (Novel, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return Novel{}, ErrInvalidInput
	}

	chapterTitle := strings.TrimSpace(in.Title)

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
		if chapterTitle == "" {
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
	return Novel{}, ErrNovelNotFound
}

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
