package server

import (
	"configs-server/data"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	dataManager *data.Manager
	router      *gin.Engine
}

func NewServer(manager *data.Manager) *Server {
	server := &Server{dataManager: manager}
	server.configureServer()
	return server
}

func (s *Server) configureServer() {
	f, _ := os.Create("config-server.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.Default()

	s.router.GET("/config/", s.handleReadConfig)
	s.router.POST("/config/", s.handleCreateConfig)
	s.router.PUT("/config/", s.handleUpdateConfig)
	s.router.DELETE("/config/", s.handleDeleteConfig)
}

func (s *Server) StartServer() error {
	return s.router.Run(fmt.Sprintf(":8080"))
}

func (s *Server) handleCreateConfig(c *gin.Context) {
	var config data.ConfigEntity
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect config format"})
		return
	}

	if err := s.dataManager.CreateConfig(config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (s *Server) handleReadConfig(c *gin.Context) {
	service := c.Query("service")
	version := c.Query("v")

	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specify 'service' parameter"})
		return
	}

	var config data.ConfigEntity
	var err error
	if version != "" {
		versionInt, err := strconv.ParseInt(version, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse version"})
		}
		config, err = s.dataManager.ReadConfigWithVersion(service, int(versionInt))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No such config"})
			return
		}
	} else {
		config, err = s.dataManager.ReadConfig(service)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No such config"})
			return
		}
	}

	c.JSON(http.StatusOK, config)
}

func (s *Server) handleUpdateConfig(c *gin.Context) {
	var config data.ConfigEntity
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect config format"})
		return
	}

	config, err := s.dataManager.UpdateConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}
func (s *Server) handleDeleteConfig(c *gin.Context) {
	service := c.Query("service")

	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specify 'service' parameter"})
		return
	}

	configs, err := s.dataManager.DeleteConfig(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}
