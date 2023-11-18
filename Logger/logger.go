package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var Zero zerolog.Logger
var Chain alice.Chain

func init() {
	appLogFile, _ := os.OpenFile(
		"app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	appLog := zerolog.MultiLevelWriter(appLogFile)
	Zero = zerolog.New(appLog).With().Timestamp().Logger()

	httpLogFile, _ := os.OpenFile(
		"http.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	httpLog := zerolog.MultiLevelWriter(httpLogFile)
	httpLogger := zerolog.New(httpLog).With().Timestamp().Logger()

	Chain = alice.New()
	Chain = Chain.Append(hlog.NewHandler(httpLogger))

	Chain = Chain.Append(
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
	)
	Chain = Chain.Append(hlog.RemoteAddrHandler("ip"))
	Chain = Chain.Append(hlog.UserAgentHandler("user_agent"))
	Chain = Chain.Append(hlog.RefererHandler("referer"))
	Chain = Chain.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

}
