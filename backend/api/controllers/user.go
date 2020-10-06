package controllers

import (
	"encoding/json"
	"github.com/resssoft/mediaArchive/auth"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/resssoft/mediaArchive/services"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserRouter struct {
	repo repositories.UserRepository
	app  services.IUserService
}

func NewUserRoute(repo repositories.UserRepository, app services.IUserService) *UserRouter {
	return &UserRouter{repo: repo, app: app}
}

func (r *UserRouter) UserInfo(ctx *fasthttp.RequestCtx) {
	userID := string(ctx.QueryArgs().Peek("userId"))
	currentUserID := ctx.UserValue("userId").(string)
	if userID == "" {
		userID = currentUserID
	}
	user, err := r.app.Get(userID, currentUserID)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 33))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, user)
}

func (r *UserRouter) AddUser(ctx *fasthttp.RequestCtx) {
	newUser := new(models.User)
	err := json.Unmarshal(ctx.PostBody(), newUser)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 31))
		return
	}
	currentUserId := ctx.UserValue("userId").(string)
	err = r.app.Add(*newUser, currentUserId)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 32))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, "OK")
}

func (r *UserRouter) Login(ctx *fasthttp.RequestCtx) {
	var userCredentials struct {
		Username string
		Password string
	}
	if len(ctx.PostBody()) > 0 {
		if err := json.Unmarshal(ctx.PostBody(), &userCredentials); err != nil {
			writeJsonResponse(ctx, http.StatusBadRequest, getError("Invalid json provided", 11))
			return
		}
	} else {
		userCredentials.Username = string(ctx.FormValue("Username"))
		userCredentials.Password = string(ctx.FormValue("Password"))
	}
	if userCredentials.Username == "" || userCredentials.Password == "" {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("Invalid arguments", 12))
		return
	}

	user, err := r.app.CheckCredentials(userCredentials.Username, userCredentials.Password)
	if err != nil {
		writeJsonResponse(ctx, http.StatusUnprocessableEntity, getError(err.Error(), 15))
		return
	}
	session := auth.NewSessionByRequest(ctx, user)
	if err := auth.SaveSession(session); err != nil {
		writeJsonResponse(ctx, http.StatusInternalServerError, getError(err.Error(), 15))
		return
	}
	expr, accessToken, err := auth.NewAccessToken(user.Role.PermsToString(), user.Id.Hex(), user.Lang, session.ID)
	if err != nil {
		writeJsonResponse(ctx, http.StatusInternalServerError, getError(err.Error(), 15))
		return
	}
	token := string(accessToken)
	writeJsonResponse(ctx, http.StatusOK, auth.RequestJwt{
		AccessToken: token,
		ExpiresIn:   expr.Unix(),
		IssuedAt:    expr.Unix(),
	})
}

func (*UserRouter) RefreshToken(ctx *fasthttp.RequestCtx) {
	sessionId, _ := ctx.UserValue("session").(string)
	session, err := auth.GetSession(sessionId)
	if err != nil || !session.IsValid(ctx) {
		ctx.SetStatusCode(http.StatusUnauthorized)
		return
	}

	expr, access, err := auth.NewAccessToken(session.Perms, session.UserId, session.UserLang, session.ID)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	if _, err := auth.UpdateSessionExpiration(session, config.JwtRtExpires()); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	token := string(access)
	writeJsonResponse(ctx, http.StatusOK, auth.RequestJwt{
		AccessToken: token,
		ExpiresIn:   expr.Unix(),
		IssuedAt:    expr.Unix(),
	})
}

func (*UserRouter) Logout(ctx *fasthttp.RequestCtx) {
	session, _ := ctx.UserValue("session").(string)
	if err := auth.DeleteSessionById(session); err != nil {
		writeJsonResponse(ctx, http.StatusInternalServerError, getError("unauthorized", 17))
		return
	}
}
