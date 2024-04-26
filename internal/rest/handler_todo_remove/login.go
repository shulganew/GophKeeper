package handler

/*
type HandlerLogin struct {
	usrSrv *services.UserService
	conf   config.Config
}

func NewHandlerLogin(conf config.Config, userServ *services.UserService) *HandlerLogin {

	return &HandlerLogin{usrSrv: userServ, conf: conf}
}

// @Summary      Login user handler
// @Description  hendler for user login
// @Tags         api
// @Success      200 success login
// @Failure      401 bad login or pass
// @Failure      404
// @Router       /{id} [get]
func (h *HandlerLogin) LoginUser(res http.ResponseWriter, req *http.Request) {
	var user entities.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		// If can't decode 400
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	userID, isValid := h.usrSrv.IsValid(req.Context(), user.Login, user.Password)
	if !isValid {
		// Wrond user login or password 401
		http.Error(res, "Wrong login or password", http.StatusUnauthorized)
		return
	}

	user.UUID = *userID

	zap.S().Debug("Login sucsess, user id is: ", userID)
	jwt, _ := services.BuildJWTString(*userID, h.conf.PassJWT)

	res.Header().Add("Content-Type", "text/plain")
	res.Header().Add("Authorization", config.AuthPrefix+jwt)

	// set status code 200
	res.WriteHeader(http.StatusOK)

	_, err := res.Write([]byte("User loged in."))
	if err != nil {
		zap.S().Errorln("Can't write to response in LoginUser handler", err)
	}
}
*/
