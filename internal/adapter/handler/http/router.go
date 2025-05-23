package httpserver

import (
	"1337bo4rd/internal/core/service"
	"net/http"
)

func NewRouter(
	postHandler PostHandler,
	userSvc *service.UserService,
) http.Handler {
	mux := http.NewServeMux()

	mw := SessionMiddleware(userSvc)
	mux.Handle("/post/", mw(http.HandlerFunc(postHandler.HandlePost)))
	mux.Handle("/archive", mw(http.HandlerFunc(postHandler.HandleArchive)))
	mux.Handle("/create", mw(http.HandlerFunc(postHandler.HandleCreate)))
	mux.Handle("/", mw(http.HandlerFunc(postHandler.HandleCatalog)))

	return mux
}
