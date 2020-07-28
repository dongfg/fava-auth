package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/markbates/pkger"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type proxy struct {
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Token   string `json:"token"`
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" && req.URL.Path == "/login" {
		handleLogin(wr, req)
		return
	}

	if req.Method == "GET" && req.URL.Path == "/login.html" {
		box, _ := pkger.Open("/public")
		http.FileServer(box).ServeHTTP(wr, req)
		return
	}

	favaAuth, _ := req.Cookie("x-auth-fava")
	if favaAuth == nil || !validateToken(favaAuth.Value) {
		http.Redirect(wr, req, "/login.html", http.StatusSeeOther)
		return
	}

	target, _ := url.Parse("http://" + config.Server)
	req.Header.Set("Host", config.Server)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(wr, req)
}

func handleLogin(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)

	var loginResp = loginResp{}

	var loginReq loginReq
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		log.Println(err)
		loginResp.Success = false
		loginResp.Msg = "Login failed, error request body."
	} else {
		if loginReq.Username == config.Username && loginReq.Password == config.Password {
			loginResp.Success = true
			loginResp.Msg = "Login success"
			loginResp.Token = generateToken()
		} else {
			loginResp.Success = false
			loginResp.Msg = "Login failed, password error."
		}
	}
	err = json.NewEncoder(wr).Encode(loginResp)
	if err != nil {
		log.Println(err)
	}
}

func startServer() {
	proxy := &proxy{}
	setupServer(&http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: proxy,
	})
}

func setupServer(srv *http.Server) {
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Start Server @ %s\n", srv.Addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("Server exiting")
}
