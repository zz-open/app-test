package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Name    string
	Version string
	Port    int
}

var (
	cfgFile = "config.yaml"
	config  Config
)

func main() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("读取配置文件失败:", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("解析失败配置文件, %v", err)
	}

	http.HandleFunc("/info", InfoHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/log", LogHandler)

	addr := fmt.Sprintf(":%d", config.Port)
	server := &http.Server{Addr: addr, Handler: nil}
	go func() {
		err := server.ListenAndServe()
		// 通过Shutdown 或者 Close 主动调用抛出
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("%s %s, cause:[%v]\n", getServiceMsg(), "http server has been close", err)
		} else {
			log.Fatalf("%s %s, cause:[%v]", getServiceMsg(), "服务启动失败", err)
		}
	}()

	log.Printf("%s listen: %s\n", getServiceMsg(), addr)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Shutdown Server...")

	//创建超时上下文，Shutdown可以让未处理的连接在这个时间内关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting...")
}

// InfoHandler
// 假设是业务接口
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"msg": fmt.Sprintf("%s %s", getServiceMsg(), "info接口"),
	}

	log.Println("Info Req Host:", r.Host)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"msg": fmt.Sprintf("%s %s", getServiceMsg(), "I am Health"),
	}

	log.Println("health Req Host:", r.Host)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	time.Sleep(time.Second * 2)

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Log Req Host:", r.Host)

	infoMap := map[string]any{
		"time": time.Now().Format(time.DateTime),
		"msg":  RandomString(10),
	}

	logData, err := json.Marshal(infoMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = _log(logData)

	data := map[string]any{
		"msg": fmt.Sprintf("%s %s", getServiceMsg(), "Log成功写入"),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func _log(content []byte) error {
	if len(content) == 0 {
		return nil
	}

	fileName := "hello.log" + time.Now().Format(time.DateOnly)
	file, err := os.OpenFile("./"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	_, _ = file.WriteString("\n")
	return nil
}

func getServiceMsg() string {
	return fmt.Sprintf("%s %s", config.Name, config.Version)
}
