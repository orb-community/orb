package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

func openLogFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		defer f.Close()
		return nil, err
	}
	return f, nil
}

type Server struct {
	config Config
	logger *log.Logger
	engine *gin.Engine
	closer func()
}

func (s *Server) Close() {
	s.closer()
}

func New(c Config) (*Server, error) {
	s := &Server{config: c}

	// setup error logger
	errorLogFile, err := openLogFile(c.Loggers.ErrorLog.Path)
	if err != nil {
		defer errorLogFile.Close()
		return nil, err
	}

	logger := log.New()
	if level, err := log.ParseLevel(c.Loggers.ErrorLog.Level); err != nil {
		return nil, err
	} else {
		logger.SetLevel(level)
	}

	logger.SetOutput(errorLogFile)
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	s.logger = logger

	// setup access logger
	accessLogFile, err := openLogFile(c.Loggers.AccessLog.Path)
	if err != nil {
		defer accessLogFile.Close()
		return nil, err
	}
	// setup gin
	gin.DefaultWriter = accessLogFile
	gin.DefaultErrorWriter = s.logger.Writer()

	r := gin.New()
	r.Use(gin.Recovery())

	s.engine = r
	s.RegisterRoutes()

	// close all the things
	s.closer = func() {
		accessLogFile.Close()
		errorLogFile.Close()
	}

	return s, nil

}

func (s *Server) Serve() error {
	c := s.config
	if c.Listener.TLSConfig.Key != "" && c.Listener.TLSConfig.Cert != "" {
		return s.engine.RunTLS(c.Listener.BindAddr, c.Listener.TLSConfig.Cert, c.Listener.TLSConfig.Key)
	}

	if c.Listener.TLSConfig.AutoTLS.Enabled {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(c.Listener.TLSConfig.AutoTLS.Domain),
			Cache:      autocert.DirCache(c.Listener.TLSConfig.AutoTLS.CertCacheDir),
		}

		if err := os.MkdirAll(c.Listener.TLSConfig.AutoTLS.CertCacheDir, 0644); err != nil {
			log.Fatal(err)
		}
		srv := &http.Server{
			Addr:      c.Listener.BindAddr,
			TLSConfig: m.TLSConfig(),
			Handler:   s.engine,
		}
		return srv.ListenAndServeTLS("", "")
	}

	return s.engine.Run(c.Listener.BindAddr)
}
