package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	handler http.Handler
	addr    string
	svr     *http.Server
}

type Option func(*Server)

func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.handler = handler
	}
}

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func New(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	s.svr = &http.Server{
		Addr:           s.addr,
		Handler:        s.handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}

	go func() {
		log.Printf("http://localhost%s\n", s.addr)

		// s.get_ip()
		err := s.svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("svr.ListenAndServe err: %v", err)
		}
	}()

}

func (s *Server) Stop(disconnects ...func() error) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// 接受 syscall.SIGINT 和 syscall.SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("")
	log.Println("Shuting down server...")

	// if err := s.client.Prisma.Disconnect(); err != nil {
	// 	log.Println("could not disconnect: ", err)
	// 	// panic(err)
	// }

	for _, disconnect := range disconnects {
		// disconnect()
		if err := disconnect(); err != nil {
			log.Println("could not disconnect: ", err)
			// panic(err)
		}
	}

	// 最大时间控制，用于通知该服务端它有 5 秒的时间处理原来的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.svr.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server extting")
}

func (s *Server) get_ip() {
	// log.Println(net.IPv4allrouter.String())
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range interfaces {
		if v.Flags&net.FlagUp == 0 || v.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := v.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if !strings.Contains(ip.IP.String(), "fe80") {
				log.Printf("http://%s%s\n", ip.IP.String(), s.addr)
			}
		}
	}
}
