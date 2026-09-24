package oidcrp

import (
	"net/http"

	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	_session_web  = "_session_web"
	_redirect_uri = "_redirect_uri"
)

type AppSessionHelper struct {
	dao *oidc_storage.AppSessionDaoPebble2
}

func NewAppSessionHelper(dao *oidc_storage.AppSessionDaoPebble2) *AppSessionHelper {
	return &AppSessionHelper{dao: dao}
}

func (p *AppSessionHelper) GetWebId(w http.ResponseWriter, req *http.Request) (*oidc_storage.AppSession, error) {
	c, err := req.Cookie(_session_web)
	if err != nil {
		newS := oidc_storage.NewAppSession()
		newS.CreationDate = timestamppb.Now()
		err = p.dao.Add(&newS)
		if err != nil {
			return nil, err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     _session_web,
			Value:    newS.WebId,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		return &newS, nil
	} else {
		webId := c.Value
		oldS, err := p.dao.FindByWebId(webId)
		if err != nil {
			return nil, err
		}
		return oldS, nil
	}
}

func (p *AppSessionHelper) RedirectUri(req *http.Request) string {
	redirectUri := req.URL.Query().Get(_redirect_uri)
	return redirectUri
}

func (p *AppSessionHelper) SaveSession(session *oidc_storage.AppSession) error {
	err := p.dao.Add(session)
	if err != nil {
		return err
	}
	return nil
}
