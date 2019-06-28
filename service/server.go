package service

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
)

type Server struct {
	Server *http.Server
	Status chan int
}

func NewServer() *Server {
	// c := config.Conf
	s := &Server{
		Status: make(chan int, 1),
	}
	return s
}

func (s *Server) StartHTTP() error {
	conf := config.Conf

	gin.SetMode(gin.DebugMode)
	r := gin.New()

	pprof.Register(r)

	s.Router(r)

	s.Server = &http.Server{
		Addr:           conf.Server.HTTP.Address + ":" + conf.Server.HTTP.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go s.Server.ListenAndServe()

	return nil
}

func Shutdown(s *http.Server, status chan int) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case sig := <-ch:
			log.Entry.Error().Str("canal get a signal", sig.String()).Msg("http shut down")
			signalShutdown(&sig, s)
			return
		case sta := <-status:
			log.Entry.Error().Int("shutdown status, sta is :", sta).Msg("shut down")
			time.Sleep(time.Second)
			return
		}
	}
}

func signalShutdown(sig *os.Signal, s *http.Server) {
	switch *sig {
	case syscall.SIGQUIT, syscall.SIGTERM /*syscall.SIGSTOP,*/, syscall.SIGINT:
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		s.Shutdown(ctx)
		return

	case syscall.SIGHUP:
	default:
		return
	}
	return
}
