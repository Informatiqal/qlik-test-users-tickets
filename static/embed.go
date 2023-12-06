package frontend

import (
	"embed"
	"io/fs"
	"net/http"

	logger "github.com/informatiqal/qlik-test-users-tickets/Logger"
)

//go:generate npm i
//go:generate npm run build
//go:embed dist
var BuildFs embed.FS

func BuildHTTPFS() http.FileSystem {
	log := logger.Zero

	build, err := fs.Sub(BuildFs, "dist")
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return http.FS(build)
}
