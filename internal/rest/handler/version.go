package handler

import (
	"net/http"

	"github.com/shulganew/GophKeeper/internal/app/config"
	"go.uber.org/zap"
)

// API REST
//
// Post "/"
//
// Get  "/"
type Version struct {
	conf *config.Config
}

// Service constructor.
func NewVersion(conf *config.Config) *Version {
	return &Version{conf: conf}
}

// GET version.
// @Summary      Get origin URL by brief (short) URL
// @Description  get version of the gkeeper
// @Tags         Version
// @Success      200
// @Router       / [get]
func (u *Version) GetVersion(res http.ResponseWriter, req *http.Request) {
	zap.S().Infof("Get request from: %+v\n", req.Header)
	// set content type
	res.Header().Add("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("1.0.0"))
}
