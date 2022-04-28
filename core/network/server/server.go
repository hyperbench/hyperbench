package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/gin-gonic/gin"
	"github.com/hyperbench/hyperbench/core/controller/worker"
	"github.com/hyperbench/hyperbench/core/network"
	"github.com/mholt/archiver/v3"
	"github.com/op/go-logging"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// Server is used to receive info from master.
type Server struct {
	logger *logging.Logger
	port   string

	nonce        int
	nLock        sync.RWMutex
	workerHandle worker.Worker
	fp           string
	index        int
}

const (
	idleNonce = -1
)

// NewServer create Server.
func NewServer(port int) *Server {
	if port == 0 {
		port = 8080
	}

	return &Server{
		logger: fcom.GetLogger("svr"),
		nonce:  idleNonce,
		port:   ":" + cast.ToString(port),
	}
}

// Start start to listen.
func (s *Server) Start() error {
	r := gin.Default()

	r.POST(network.SetNoncePath, func(c *gin.Context) {

		s.nLock.Lock()
		defer s.nLock.Unlock()

		if s.nonce != idleNonce {
			s.logger.Error("busy")
			c.String(http.StatusNotAcceptable, "busy")
			return
		}
		sNonce, exist := c.GetPostForm("nonce")
		if !exist {
			s.logger.Error("need nonce")
			c.String(http.StatusNotAcceptable, "need nonce")
			return
		}
		i, err := strconv.Atoi(sNonce)
		if err != nil {
			s.logger.Error("nonce error")
			c.String(http.StatusNotAcceptable, "nonce error")
			return
		}
		s.nonce = i
	})

	r.POST(network.UploadPath, func(c *gin.Context) {

		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}
		dir, _ := c.GetPostForm(network.ConfigPath)
		f, err := c.FormFile(network.FileName)
		if err != nil {
			s.logger.Error("need file")
			c.String(http.StatusNotAcceptable, "need file")
			return
		}
		s.fp = f.Filename[0 : strings.LastIndex(f.Filename, "/")+1]
		// clear benchmark path for upload
		s.removeBenchmark(s.fp)
		s.createBenchmark(s.fp)
		s.logger.Noticef("upload %v", f.Filename)
		err = c.SaveUploadedFile(f, f.Filename)
		if err != nil {
			s.logger.Error("can not save file")
			c.String(http.StatusNotAcceptable, "can not save file")
			return
		}

		s.logger.Notice("fp", f.Filename)
		err = archiver.Unarchive(f.Filename, s.fp)
		if err != nil {
			if strings.Contains(err.Error(), "file already exists") {
				s.logger.Errorf("can not open file: %v", err)
			} else {
				s.logger.Errorf("can not open file: %v", err)
				c.String(http.StatusNotAcceptable, "can not open file")
				return
			}
		}
		s.removeBenchmark(f.Filename)
		viper.AddConfigPath(s.fp)
		if dir != "" {
			viper.SetConfigFile(dir)
		}
		err = viper.ReadInConfig()
		if err != nil {
			s.logger.Errorf("can not read in config")
			c.String(http.StatusNotAcceptable, "can not read in config")
			return
		}

		c.String(http.StatusOK, "ok")
	})

	r.POST(network.InitPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}

		sIndex, exist := c.GetPostForm("index")
		if !exist {
			s.logger.Error("need index")
			c.String(http.StatusNotAcceptable, "need index")
			return
		}

		var err error
		s.index, err = strconv.Atoi(sIndex)
		if err != nil {
			s.logger.Error("invalid index")
			c.String(http.StatusNotAcceptable, "invalid index")
			return
		}

		l := len(viper.GetStringSlice(fcom.EngineURLsPath))
		if l < s.index || l == 0 {
			s.logger.Error("config error")
			c.String(http.StatusNotAcceptable, "config error")
			return
		}

		s.workerHandle, err = worker.NewLocalWorker(worker.LocalWorkerConfig{
			Index:    int64(s.index),
			Instant:  int64(viper.GetInt(fcom.EngineInstantPath) / l),
			Wait:     viper.GetDuration(fcom.EngineWaitPath),
			Cap:      int64(viper.GetInt(fcom.EngineCapPath) / l),
			Rate:     int64(viper.GetInt(fcom.EngineRatePath) / l),
			Duration: viper.GetDuration(fcom.EngineDurationPath),
		})

		if err != nil {
			s.logger.Error("create worker error: %v", err)
			c.String(http.StatusNotAcceptable, "create worker error: %v", err)
			return
		}
	})

	r.POST(network.SetContextPath, func(c *gin.Context) {

		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}

		if s.workerHandle == nil {
			s.logger.Error("worker is not exist")
			c.String(http.StatusUnauthorized, "worker is not exist")
			return
		}

		ctx, exist := c.GetPostForm("context")
		if !exist {
			s.logger.Error("need context")
			c.String(http.StatusNotAcceptable, "need context")
			return
		}

		err := s.workerHandle.SetContext(network.Hex2Bytes(ctx))
		if err != nil {
			s.logger.Error("set context error")
			c.String(http.StatusNotAcceptable, "set error: %v", err)
			return
		}
		c.String(http.StatusOK, "ok")
	})

	r.POST(network.BeforeRunPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}
		if s.workerHandle == nil {
			s.logger.Error("worker is not exist")
			c.String(http.StatusUnauthorized, "worker is not exist")
			return
		}
		// nolint
		go s.workerHandle.BeforeRun()

		c.String(http.StatusOK, "ok")
	})

	r.POST(network.DoPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}
		if s.workerHandle == nil {
			s.logger.Error("worker is not exist")
			c.String(http.StatusUnauthorized, "worker is not exist")
			return
		}
		// nolint
		go s.workerHandle.Do()

		c.String(http.StatusOK, "ok")
	})

	r.POST(network.StatisticsPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}
		if s.workerHandle == nil {
			s.logger.Error("worker is not exist")
			c.String(http.StatusUnauthorized, "worker is not exist")
			return
		}

		sent, missed := s.workerHandle.Statistics()
		Sent := strconv.FormatInt(sent, 10)
		Missed := strconv.FormatInt(missed, 10)
		dm := map[string]interface{}{
			"sent":   Sent,
			"missed": Missed,
		}

		c.JSON(http.StatusOK, dm)
	})

	r.POST(network.AfterRunPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}
		if s.workerHandle == nil {
			s.logger.Error("worker is not exist")
			c.String(http.StatusUnauthorized, "worker is not exist")
			return
		}
		// nolint
		go s.workerHandle.AfterRun()

		c.String(http.StatusOK, "ok")
	})

	r.POST(network.CheckoutCollectorPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}
		if s.workerHandle == nil {
			s.logger.Error("worker is not exist")
			c.String(http.StatusUnauthorized, "worker is not exist")
			return
		}

		col, valid, err := s.workerHandle.CheckoutCollector()
		if err != nil {
			//todo add something for err
		}
		var t, data string
		if col != nil {
			t = col.Type()
			data = network.Bytes2Hex(col.Serialize())
		}

		dm := map[string]interface{}{
			"type":  t,
			"col":   data,
			"valid": valid,
		}

		c.JSON(http.StatusOK, dm)
	})

	r.POST(network.TeardownPath, func(c *gin.Context) {
		if !s.checkNonce(c) {
			s.logger.Error("busy")
			c.String(http.StatusUnauthorized, "busy")
			return
		}

		if s.workerHandle != nil {
			s.workerHandle.Teardown()
			s.workerHandle = nil
		}

		s.nonce = idleNonce
		//s.removeBenchmark()
		viper.Reset()
		c.String(http.StatusOK, "ok")
	})

	return r.Run(s.port)
}

func (s *Server) checkNonce(c *gin.Context) bool {
	if s.nonce == idleNonce {
		return false
	}
	sNonce, _ := c.GetPostForm("nonce")
	i, _ := strconv.Atoi(sNonce)
	return i == s.nonce
}

func (s *Server) createBenchmark(filename string) {
	_, err := os.Stat(filename)
	if err != nil && !os.IsExist(err) {
		_ = os.MkdirAll(filename, 0777)
	}
}

func (s *Server) removeBenchmark(fileName string) {
	_ = os.RemoveAll(fileName)
}
