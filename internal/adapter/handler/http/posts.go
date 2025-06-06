package httpserver

import (
	"1337bo4rd/internal/adapter/storage/minio"
	"1337bo4rd/internal/core/domain"
	"1337bo4rd/internal/core/port"
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	svc     port.PostService
	tmpl    *template.Template
	storage *minio.MinioClient
}

var (
	msg        = "Failed to load threads."
	statusCode = http.StatusInternalServerError
)

func NewPostHandler(svc port.PostService, mio *minio.MinioClient) *PostHandler {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	return &PostHandler{
		svc,
		tmpl,
		mio,
	}
}

func (h *PostHandler) HandleCatalog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		statusCode = http.StatusNotFound
		msg = "No such url path >,<"
		renderError(w, h.tmpl, statusCode, msg)
		return
	}
	posts, err := h.svc.ListActive()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		renderError(w, h.tmpl, statusCode, msg)
		return
	}
	h.tmpl.ExecuteTemplate(w, "catalog.html", struct{ Posts []domain.Post }{posts})
}

func (h *PostHandler) HandleArchive(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.ListPosts()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		renderError(w, h.tmpl, statusCode, msg)
		return
	}
	h.tmpl.ExecuteTemplate(w, "archive.html", struct{ Posts []domain.Post }{posts})
}

func (h *PostHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	// add logic for r.Method(POST) Addd comments!
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/post/"), "/")
	id := parts[0]

	if len(parts) > 1 && parts[1] == "comment" && r.Method == http.MethodPost {
		h.addComment(w, r, id)
		return
	}

	postWithComments, err := h.svc.GetPostWithCommentsById(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			msg = port.ErrNoPosts.Error()
			statusCode = http.StatusNotFound
		case port.ErrInvalidPostId:
			statusCode = http.StatusBadRequest
			msg = err.Error()
		}
		renderError(w, h.tmpl, statusCode, msg)
		return
	}

	// Build comment tree
	nodes := make(map[uint64]*domain.CommentNode)
	for _, c := range postWithComments.Comments {
		comment := c // copy to avoid pointer issue
		nodes[c.ID] = &domain.CommentNode{Comment: &comment}
	}

	var roots []*domain.CommentNode
	for _, node := range nodes {
		if node.ParentCommentID != 0 {
			if parent, ok := nodes[node.ParentCommentID]; ok {
				parent.Replies = append(parent.Replies, node)
				continue
			}
		}
		roots = append(roots, node)
	}

	// Template data
	data := struct {
		Post        *domain.Post
		CommentTree []*domain.CommentNode
		User        *domain.User
	}{
		Post:        &postWithComments.Post,
		CommentTree: roots,
		User:        getSession(r),
	}
	if !postWithComments.Post.ArchivedAt.IsZero() {
		h.tmpl.ExecuteTemplate(w, "archive-post.html", data)
	} else {
		h.tmpl.ExecuteTemplate(w, "post.html", data)
	}
}

func (h *PostHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.tmpl.ExecuteTemplate(w, "create-post.html", nil)
		return
	}

	userSession := getSession(r)
	title := r.FormValue("title")
	content := r.FormValue("content")
	if strings.TrimSpace(content) == "" {
		renderError(w, h.tmpl, http.StatusBadRequest, "Content cannot be empty or whitespace")
		return
	}
	if strings.TrimSpace(title) == "" {
		renderError(w, h.tmpl, http.StatusBadRequest, "Title cannot be empty or whitespace")
		return
	}

	// Minio usage
	var imageURL string
	if f, _, err := r.FormFile("image"); err == nil {
		defer f.Close()

		url, err := h.storage.UploadImage(r.Context(), f, "post-image.jpg", "image/jpeg")
		if err != nil {
			renderError(w, h.tmpl, http.StatusInternalServerError, "Error uploading image")
			return
		}

		imageURL = url
	} else if err != http.ErrMissingFile {
		renderError(w, h.tmpl, http.StatusBadRequest, "Error extracting image")
		return
	}

	post := &domain.Post{
		Title:      title,
		Content:    content,
		UserName:   userSession.Name,
		UserAvatar: userSession.Avatar,
		Image:      imageURL,
	}

	err := h.svc.CreatePost(post, userSession.ID)
	if err != nil {
		renderError(w, h.tmpl, statusCode, msg)
		return
	}
	slog.Info("created post")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) addComment(w http.ResponseWriter, r *http.Request, postID string) {
	userSession := getSession(r)
	parentStr := r.FormValue("parent_id")
	var parent uint64
	if parentStr != "" {
		if v, err := strconv.ParseUint(parentStr, 10, 64); err == nil {
			parent = v
		}
	}

	content := r.FormValue("content")

	if strings.TrimSpace(content) == "" {
		renderError(w, h.tmpl, http.StatusBadRequest, "Content cannot be empty or whitespace")
		return
	}
	comment := &domain.Comment{
		UserName:        userSession.Name,
		UserAvatar:      userSession.Avatar,
		ParentCommentID: parent,
		Content:         content,
	}

	err := h.svc.CreateComment(comment, postID, userSession.ID)
	if err != nil {
		if errors.Is(err, port.ErrInvalidPostId) {
			statusCode = http.StatusBadRequest
			msg = err.Error()
		}
		renderError(w, h.tmpl, statusCode, msg)
		return
	}
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}
