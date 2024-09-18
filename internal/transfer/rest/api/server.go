package api

import (
	"context"
	"errors"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/repository"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
)

type Server struct {
	logger         *slog.Logger
	validator      *validator.Validate
	httpServer     *http.Server
	router         *mux.Router
	redisService   repository.ScooterRepository
	rentalService  rental.RentalService
	trackerService tracker.Service
	eligibleUsers  map[string]bool
}

func NewServer(
	logger *slog.Logger,
	validator *validator.Validate,
	server *http.Server,
	router *mux.Router,
	redis repository.ScooterRepository,
	rental rental.RentalService,
	tracker tracker.Service,
	users map[string]bool,
) *Server {

	s := &Server{
		logger:         logger,
		validator:      validator,
		httpServer:     server,
		router:         router,
		redisService:   redis,
		rentalService:  rental,
		trackerService: tracker,
		eligibleUsers:  users,
	}

	s.registerRoutes()

	return s
}

// Run starts the work of the service.
func (s *Server) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-signalChan
		s.logger.Info(sig.String())

		cancel()
	}()

	waitGroup.Add(1)

	go func(running *sync.WaitGroup) {
		defer running.Done()
		s.logger.Info("Starting HTTP server.", slog.String("address", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancel()

			s.logger.Error("can't close http server", err)

			return
		}

		s.logger.Info("Stopping HTTP server.")
	}(&waitGroup)

	<-ctx.Done()

	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Error("can't shutdown gracefully", err)

		return
	}

	waitGroup.Wait()
}
