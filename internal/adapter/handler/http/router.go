package httpserver

import (
	"1337bo4rd/internal/adapter/storage/minio"
	"1337bo4rd/internal/core/service"
	"net/http"
)

func NewRouter(
	postHandler PostHandler,
	userSvc *service.UserService,
	userHandler UserHandler,
) http.Handler {
	mux := http.NewServeMux()

	mw := SessionMiddleware(userSvc)
	// post routers
	mux.Handle("/post/", mw(http.HandlerFunc(postHandler.HandlePost)))
	mux.Handle("/archive", mw(http.HandlerFunc(postHandler.HandleArchive)))
	mux.Handle("/create", mw(http.HandlerFunc(postHandler.HandleCreate)))
	mux.HandleFunc("GET /images/posts/{filename}", minio.ServePostImageHandler(postHandler.storage))
	mux.Handle("/", mw(http.HandlerFunc(postHandler.HandleCatalog)))
	// user routers
	mux.Handle("/profile", mw(http.HandlerFunc(userHandler.HandleProfile)))
	return mux
}
