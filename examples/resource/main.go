// Package main demonstrates the RESTful resource builder pattern.
package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/kolosys/helix"
	"github.com/kolosys/helix/middleware"
)

// Article represents a blog article.
type Article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Published bool   `json:"published"`
}

// Comment represents a comment on an article.
type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	Author    string `json:"author"`
	Body      string `json:"body"`
}

// In-memory storage
var (
	articles = map[int]Article{
		1: {ID: 1, Title: "Getting Started with Helix", Content: "Learn the basics...", Author: "Alice", Published: true},
		2: {ID: 2, Title: "Advanced Routing", Content: "Deep dive into routing...", Author: "Bob", Published: true},
	}
	comments = map[int][]Comment{
		1: {
			{ID: 1, ArticleID: 1, Author: "Charlie", Body: "Great article!"},
			{ID: 2, ArticleID: 1, Author: "Diana", Body: "Very helpful, thanks!"},
		},
	}
	articleMu sync.RWMutex
	nextArtID = 3
	nextComID = 3
)

func main() {
	s := helix.Default(&helix.Options{
		Addr: ":8080",
	})

	// Root endpoint
	s.GET("/", helix.HandleCtx(func(c *helix.Ctx) error {
		return c.OK(map[string]string{
			"message": "Blog API - Use /articles for articles",
		})
	}))

	// Define a resource for articles using the fluent builder
	s.Resource("/articles").
		List(listArticles).                                 // GET /articles
		Create(createArticle).                              // POST /articles
		Get(getArticle).                                    // GET /articles/{id}
		Update(updateArticle).                              // PUT /articles/{id}
		Delete(deleteArticle).                              // DELETE /articles/{id}
		Custom("POST", "/{id}/publish", publishArticle).    // POST /articles/{id}/publish
		Custom("POST", "/{id}/unpublish", unpublishArticle) // POST /articles/{id}/unpublish

	// Define a resource for comments within the API group
	api := s.Group("/api/v1")

	// Nested resource: comments under articles
	api.Resource("/articles/{articleId}/comments").
		List(listComments).
		Create(createComment).
		Get(getComment).
		Delete(deleteComment)

	// Admin-only resource with middleware - no casting needed!
	admin := s.Group("/admin", middleware.BasicAuth("admin", "secret"))

	admin.Resource("/articles").
		List(listArticles).
		Delete(deleteArticle)

	log.Println("Server starting on :8080")
	log.Println("Routes:")
	s.PrintRoutes(log.Writer())

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

// Article handlers
func listArticles(w http.ResponseWriter, r *http.Request) {
	articleMu.RLock()
	defer articleMu.RUnlock()

	list := make([]Article, 0, len(articles))
	for _, a := range articles {
		list = append(list, a)
	}

	helix.OK(w, map[string]any{
		"articles": list,
		"total":    len(list),
	})
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	id, err := helix.ParamInt(r, "id")
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid id"))
		return
	}

	articleMu.RLock()
	article, ok := articles[id]
	articleMu.RUnlock()

	if !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", id))
		return
	}

	helix.OK(w, article)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	req, err := helix.BindJSON[struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}](r)
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid request body"))
		return
	}

	articleMu.Lock()
	article := Article{
		ID:        nextArtID,
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		Published: false,
	}
	articles[article.ID] = article
	nextArtID++
	articleMu.Unlock()

	helix.Created(w, article)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	id, err := helix.ParamInt(r, "id")
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid id"))
		return
	}

	req, err := helix.BindJSON[struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}](r)
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid request body"))
		return
	}

	articleMu.Lock()
	defer articleMu.Unlock()

	article, ok := articles[id]
	if !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", id))
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	articles[id] = article

	helix.OK(w, article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := helix.ParamInt(r, "id")
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid id"))
		return
	}

	articleMu.Lock()
	defer articleMu.Unlock()

	if _, ok := articles[id]; !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", id))
		return
	}

	delete(articles, id)
	helix.NoContent(w)
}

func publishArticle(w http.ResponseWriter, r *http.Request) {
	id, err := helix.ParamInt(r, "id")
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid id"))
		return
	}

	articleMu.Lock()
	defer articleMu.Unlock()

	article, ok := articles[id]
	if !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", id))
		return
	}

	article.Published = true
	articles[id] = article

	helix.OK(w, article)
}

func unpublishArticle(w http.ResponseWriter, r *http.Request) {
	id, err := helix.ParamInt(r, "id")
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid id"))
		return
	}

	articleMu.Lock()
	defer articleMu.Unlock()

	article, ok := articles[id]
	if !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", id))
		return
	}

	article.Published = false
	articles[id] = article

	helix.OK(w, article)
}

// Comment handlers
func listComments(w http.ResponseWriter, r *http.Request) {
	articleID, _ := strconv.Atoi(helix.Param(r, "articleId"))

	articleMu.RLock()
	defer articleMu.RUnlock()

	if _, ok := articles[articleID]; !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", articleID))
		return
	}

	list := comments[articleID]
	if list == nil {
		list = []Comment{}
	}

	helix.OK(w, map[string]any{
		"comments": list,
		"total":    len(list),
	})
}

func getComment(w http.ResponseWriter, r *http.Request) {
	articleID, _ := strconv.Atoi(helix.Param(r, "articleId"))
	commentID, _ := strconv.Atoi(helix.Param(r, "id"))

	articleMu.RLock()
	defer articleMu.RUnlock()

	list := comments[articleID]
	for _, c := range list {
		if c.ID == commentID {
			helix.OK(w, c)
			return
		}
	}

	helix.WriteProblem(w, helix.NotFoundf("comment %d not found", commentID))
}

func createComment(w http.ResponseWriter, r *http.Request) {
	articleID, _ := strconv.Atoi(helix.Param(r, "articleId"))

	req, err := helix.BindJSON[struct {
		Author string `json:"author"`
		Body   string `json:"body"`
	}](r)
	if err != nil {
		helix.WriteProblem(w, helix.BadRequestf("invalid request body"))
		return
	}

	articleMu.Lock()
	defer articleMu.Unlock()

	if _, ok := articles[articleID]; !ok {
		helix.WriteProblem(w, helix.NotFoundf("article %d not found", articleID))
		return
	}

	comment := Comment{
		ID:        nextComID,
		ArticleID: articleID,
		Author:    req.Author,
		Body:      req.Body,
	}
	comments[articleID] = append(comments[articleID], comment)
	nextComID++

	helix.Created(w, comment)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	articleID, _ := strconv.Atoi(helix.Param(r, "articleId"))
	commentID, _ := strconv.Atoi(helix.Param(r, "id"))

	articleMu.Lock()
	defer articleMu.Unlock()

	list := comments[articleID]
	for i, c := range list {
		if c.ID == commentID {
			comments[articleID] = append(list[:i], list[i+1:]...)
			helix.NoContent(w)
			return
		}
	}

	helix.WriteProblem(w, helix.NotFoundf("comment %d not found", commentID))
}
