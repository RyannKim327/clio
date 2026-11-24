package site

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreLifecycle(t *testing.T) {
	t.Setenv("CLIO_DATA_FILE", filepath.Join(t.TempDir(), "novels.json"))
	if err := writeDatabase(novelDatabase{}); err != nil {
		t.Fatalf("seed database: %v", err)
	}

	novel, err := CreateNovel(NovelInput{
		Title:   "A Tale",
		Author:  "Someone",
		Content: "First paragraph.\n\nSecond paragraph.",
	})
	if err != nil {
		t.Fatalf("create novel: %v", err)
	}
	if !strings.HasPrefix(novel.ID, "a-tale-") {
		t.Fatalf("unexpected id: %s", novel.ID)
	}
	if len(novel.Chapters) != 1 || novel.Chapters[0].Title != "Chapter 1" {
		t.Fatalf("unexpected first chapter: %+v", novel.Chapters)
	}
	if len(novel.Chapters[0].Paragraphs) != 2 {
		t.Fatalf("expected two paragraphs, got %d", len(novel.Chapters[0].Paragraphs))
	}

	novels, err := ListNovels()
	if err != nil || len(novels) != 1 {
		t.Fatalf("list novels: %v (%d)", err, len(novels))
	}

	novel, err = AddChapter(novel.ID, ChapterInput{Content: "More story."})
	if err != nil {
		t.Fatalf("add chapter: %v", err)
	}
	if len(novel.Chapters) != 2 || novel.Chapters[1].Title != "Chapter 2" {
		t.Fatalf("unexpected chapters: %+v", novel.Chapters)
	}

	novel, err = UpdateChapter(novel.ID, 1, ChapterInput{Title: "Beginning", Content: "Rewritten."})
	if err != nil {
		t.Fatalf("update chapter: %v", err)
	}
	if novel.Chapters[0].Title != "Beginning" {
		t.Fatalf("chapter title was not updated: %+v", novel.Chapters[0])
	}

	stored, err := GetNovel(novel.ID)
	if err != nil || stored.Chapters[0].Title != "Beginning" {
		t.Fatalf("persisted novel: %v (%+v)", err, stored)
	}
	if chapterText(stored.Chapters[0]) != "Rewritten." {
		t.Fatalf("unexpected chapter text: %q", chapterText(stored.Chapters[0]))
	}

	if _, err := UpdateChapter(novel.ID, 9, ChapterInput{Content: "Nope"}); !errors.Is(err, ErrChapterNotFound) {
		t.Fatalf("expected chapter not found, got %v", err)
	}
	if _, err := CreateNovel(NovelInput{Title: " ", Content: "text"}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
	if _, err := GetNovel("missing"); !errors.Is(err, ErrNovelNotFound) {
		t.Fatalf("expected novel not found, got %v", err)
	}
}
