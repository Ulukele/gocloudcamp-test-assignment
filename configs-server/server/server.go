package server

import (
	"configs-server/data"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	f, _ := os.Create("/logs/config-server.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.Default()

	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.GET("/config/", s.handleReadConfig)
	s.router.POST("/config/", s.handleCreateConfig)
	s.router.PUT("/config/", s.handleUpdateConfig)
	s.router.DELETE("/config/", s.handleDeleteConfig)
}

func (s *Server) StartServer(addr string) error {
	return s.router.Run(addr)
}

// CreateConfig             godoc
// @Summary      Create config
// @Description  Responds with the config as JSON.
// @Tags         configs
// @Produce      json
// @Param        config  body      models.ConfigDef  true  "Config JSON"
// @Success      200 {object} models.Config
// @Router       /config/ [post]
func (s *Server) handleCreateConfig(c *gin.Context) {
	var config data.ConfigEntity
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect config format"})
		return
	}

	if err := s.dataManager.CreateConfig(config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Such config already exists"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// ReadConfig             godoc
// @Summary      Reads the latest version of the config
// @Description  Responds with the config as JSON.
// @Tags         configs
// @Produce      json
// @Param        service  path      string  true  "search config by service"
// @Success      200 {object} models.Config
// @Router       /config/ [get]
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "There is no such config"})
			return
		}
	} else {
		config, err = s.dataManager.ReadConfig(service)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There is no such config"})
			return
		}
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfig             godoc
// @Summary      Creates new config version
// @Description  Responds with the config as JSON.
// @Tags         configs
// @Produce      json
// @Param        config  body      models.ConfigDef  true  "Config JSON"
// @Success      200 {object} models.Config
// @Router       /config/ [put]
func (s *Server) handleUpdateConfig(c *gin.Context) {
	var config data.ConfigEntity
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect config format"})
		return
	}

	config, err := s.dataManager.UpdateConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update such config"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteConfig             godoc
// @Summary      Deletes all config versions
// @Description  Responds with the all configs versions as JSON.
// @Tags         configs
// @Produce      json
// @Param        service  path      string  true  "delete configs by service"
// @Success      200 {array} models.Config
// @Router       /config/ [delete]
func (s *Server) handleDeleteConfig(c *gin.Context) {
	service := c.Query("service")

	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specify 'service' parameter"})
		return
	}

	configs, err := s.dataManager.DeleteConfig(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "There is no such config"})
		return
	}

	c.JSON(http.StatusOK, configs)
}
