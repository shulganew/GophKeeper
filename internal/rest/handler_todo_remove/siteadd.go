package handler

/*
import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/GophKeeper/internal/app/config"
	"github.com/shulganew/GophKeeper/internal/entities"
	"github.com/shulganew/GophKeeper/internal/services"
	"go.uber.org/zap"
)

type HandlerSiteAdd struct {
	siteSrv *services.SiteService
	conf    config.Config
}

func NewSiteAdd(conf config.Config, siteSrv *services.SiteService) *HandlerSiteAdd {
	return &HandlerSiteAdd{siteSrv: siteSrv, conf: conf}
}

// @Summary      Add site credentials
// @Description  Add credentioals for siet - login and password
// @Tags         api
// @Success      201 Creaded - addted for user
// @Failure      401 User not auth
// @Failure      400 JSON corrupted.

// @Router       / [get]
func (s *HandlerSiteAdd) SiteAdd(res http.ResponseWriter, req *http.Request) {
	// Check registration.
	userID, isRegistered := services.CheckUserAuth(req.Context())
	if isRegistered {
		http.Error(res, "JWT not found.", http.StatusUnauthorized)
		return
	}

	var site entities.Site

	if err := json.NewDecoder(req.Body).Decode(&site); err != nil {
		// If can't decode 400
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.siteSrv.AddSite(req.Context(), userID, site.SiteURL, site.SLogin, site.SPw)
	if err != nil {

		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "text/plain")

	// set status code 200
	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte("User loged in."))
	if err != nil {
		zap.S().Errorln("Can't write to response in LoginUser handler", err)
	}
}
*/
