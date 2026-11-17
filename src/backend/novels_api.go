package backend

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RyannKim327/constellation/config/middleware"
	"github.com/RyannKim327/constellation/config/utils"

	"github.com/gin-gonic/gin"
)

var NovelsAPI = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "api/novels",
	Permission: utils.PERMISSION_ALL,
	Handler:    listNovels,
}

var NovelAPI = middleware.Route{
	Method:     utils.METHOD_GET,
	Path:       "api/novels/:id",
	Permission: utils.PERMISSION_ALL,
	Handler:    getNovel,
}

var CreateNovelAPI = middleware.Route{
	Method:     utils.METHOD_POST,
	Path:       "api/novels",
	Permission: utils.PERMISSION_ALL,
	Handler:    createNovel,
}

var AddChapterAPI = middleware.Route{
	Method:     utils.METHOD_POST,
	Path:       "api/novels/:id/chapters",
	Permission: utils.PERMISSION_ALL,
	Handler:    addChapter,
}

var UpdateChapterAPI = middleware.Route{
	Method:     utils.METHOD_PUT,
	Path:       "api/novels/:id/chapters/:chapter",
	Permission: utils.PERMISSION_ALL,
	Handler:    updateChapter,
}

// UpdateChapterFormAPI lets plain HTML forms save an edited chapter.
var UpdateChapterFormAPI = middleware.Route{
	Method:     utils.METHOD_POST,
	Path:       "api/novels/:id/chapters/:chapter/update",
	Permission: utils.PERMISSION_ALL,
	Handler:    updateChapter,
}

func wantsJSON(ctx *gin.Context) bool {
	return strings.Contains(ctx.GetHeader("Content-Type"), "application/json")
}

func listNovels(ctx *gin.Context) (string, utils.PARAMS) {
	novels, err := ListNovels()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", nil
	}
	ctx.JSON(http.StatusOK, gin.H{"novels": novels})
	return "", nil
}

func getNovel(ctx *gin.Context) (string, utils.PARAMS) {
	novel, err := GetNovel(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, ErrNovelNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "novel not found"})
			return "", nil
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", nil
	}
	ctx.JSON(http.StatusOK, gin.H{"novel": novel})
	return "", nil
}

func createNovel(ctx *gin.Context) (string, utils.PARAMS) {
	var input NovelInput
	if wantsJSON(ctx) {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return "", nil
		}
	} else {
		input = NovelInput{
			Title:        ctx.PostForm("title"),
			Author:       ctx.PostForm("author"),
			Description:  ctx.PostForm("description"),
			ChapterTitle: ctx.PostForm("chapter_title"),
			Content:      ctx.PostForm("content"),
		}
	}

	novel, err := CreateNovel(input)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidInput) {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return "", nil
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusCreated, gin.H{"novel": novel})
		return "", nil
	}
	ctx.Redirect(http.StatusSeeOther, "/read/"+novel.ID)
	return "", nil
}

func addChapter(ctx *gin.Context) (string, utils.PARAMS) {
	var input ChapterInput
	if wantsJSON(ctx) {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return "", nil
		}
	} else {
		input = ChapterInput{
			Title:   ctx.PostForm("chapter_title"),
			Content: ctx.PostForm("content"),
		}
	}

	novel, err := AddChapter(ctx.Param("id"), input)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, ErrNovelNotFound):
			status = http.StatusNotFound
		case errors.Is(err, ErrInvalidInput):
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return "", nil
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusCreated, gin.H{"novel": novel})
		return "", nil
	}
	ctx.Redirect(http.StatusSeeOther, "/read/"+novel.ID+"?chapter="+strconv.Itoa(len(novel.Chapters)))
	return "", nil
}

func updateChapter(ctx *gin.Context) (string, utils.PARAMS) {
	chapterNumber, err := strconv.Atoi(ctx.Param("chapter"))
	if err != nil || chapterNumber < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter number"})
		return "", nil
	}

	var input ChapterInput
	if wantsJSON(ctx) {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return "", nil
		}
	} else {
		input = ChapterInput{
			Title:   ctx.PostForm("chapter_title"),
			Content: ctx.PostForm("content"),
		}
	}

	novel, err := UpdateChapter(ctx.Param("id"), chapterNumber, input)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, ErrNovelNotFound), errors.Is(err, ErrChapterNotFound):
			status = http.StatusNotFound
		case errors.Is(err, ErrInvalidInput):
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return "", nil
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"novel": novel})
		return "", nil
	}
	ctx.Redirect(http.StatusSeeOther, "/read/"+novel.ID+"?chapter="+strconv.Itoa(chapterNumber))
	return "", nil
}
