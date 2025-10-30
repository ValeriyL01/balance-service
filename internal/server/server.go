package server

import (
	"fmt"
	"net/http"

	"github.com/ValeriyL01/balance-service/internal/handlers"
)

type Server struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewServer(port string, handler *handlers.Handler, userHandler *handlers.UserHandler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/balance", handler.GetUserBalance)
	mux.HandleFunc("/deposit", handler.DepositBalance)
	mux.HandleFunc("/withdraw", handler.WithdrawBalance)
	mux.HandleFunc("/transfer", handler.TransferMoney)
	mux.HandleFunc("/transactions", handler.GetTransactionUser)

	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	return &Server{mux: mux, server: &server}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()

}
