package httpserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/docktermj/demo-entity-search/entitysearchservice"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// HttpServerImpl is the default implementation of the HttpServer interface.
type HttpServerImpl struct {
	AllowedHostnames     []string
	Arguments            []string
	Command              string
	ConnectionErrorLimit int
	EnableAll            bool
	EnableEntitySearch   bool
	HtmlTitle            string
	KeepalivePingTimeout int
	MaxBufferSizeBytes   int
	ServerAddress        string
	ServerPort           int
	TtyOnly              bool
	UrlRoutePrefix       string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *HttpServerImpl) getEntitySearchMux(ctx context.Context) *http.ServeMux {
	service := &entitysearchservice.HttpServiceImpl{}
	return service.Handler(ctx)
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method serves the httpservice over HTTP.

Input
  - ctx: A context to control lifecycle.
*/

func (httpServer *HttpServerImpl) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()

	userMessage := ""

	// Enable EntitySearch.

	// if httpServer.EnableAll || httpServer.EnableEntitySearch {
	// entitySearchMux := httpServer.getEntitySearchMux(ctx)
	// rootMux.Handle("/entity-search/", http.StripPrefix("/entity-search", entitySearchMux))
	// userMessage = fmt.Sprintf("%sServing EntitySearch at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search")
	// }

	entitySearchMux := httpServer.getEntitySearchMux(ctx)
	rootMux.Handle("/entity-search/", http.StripPrefix("/entity-search", entitySearchMux))
	userMessage = fmt.Sprintf("%sServing EntitySearch at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, "entity-search")

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}
	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessage = fmt.Sprintf("%sStarting server on interface:port '%s'...\n", userMessage, listenOnAddress)
	fmt.Println(userMessage)
	server := http.Server{
		Addr:    listenOnAddress,
		Handler: addIncomingRequestLogging(rootMux),
	}

	// Start a web browser.  Unless disabled.

	// if !httpServer.TtyOnly {
	// 	_ = browser.OpenURL(fmt.Sprintf("http://localhost:%d", httpServer.ServerPort))
	// }

	return server.ListenAndServe()
}
