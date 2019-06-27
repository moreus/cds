package api

import (
	"context"
	"net/http"

	"github.com/ovh/cds/engine/api/authentication"
	"github.com/ovh/cds/engine/api/user"
	"github.com/ovh/cds/sdk/log"

	"github.com/ovh/cds/sdk"

	"github.com/ovh/cds/engine/service"
)

func (api *API) postAuthBuiltinSigninHandler() service.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// Get the consumer builtin driver
		driver, ok := api.AuthenticationDrivers[sdk.ConsumerBuiltin]
		if !ok {
			return sdk.WithStack(sdk.ErrNotFound)
		}

		// Extract and validate signin request
		var req sdk.AuthConsumerSigninRequest
		if err := service.UnmarshalBody(r, &req); err != nil {
			return err
		}
		if err := driver.CheckSigninRequest(req); err != nil {
			return err
		}
		// Convert code to external user info
		userInfo, err := driver.GetUserInfo(req)
		if err != nil {
			return err
		}

		tx, err := api.mustDB().Begin()
		if err != nil {
			return sdk.WithStack(err)
		}
		defer tx.Rollback() // nolint

		// Check if a consumer exists for consumer type and external user identifier
		log.Debug("loading consumer from %+v", userInfo)
		consumer, err := authentication.LoadConsumerByID(ctx, tx, userInfo.ExternalID)
		if err != nil {
			return err
		}

		// Generate a new session for consumer
		session, err := authentication.NewSession(tx, consumer, driver.GetSessionDuration())
		if err != nil {
			return err
		}

		// Generate a jwt for current session
		jwt, err := authentication.NewSessionJWT(session)
		if err != nil {
			return err
		}

		// Set a cookie with the jwt token
		http.SetCookie(w, &http.Cookie{
			Name:    jwtCookieName,
			Value:   jwt,
			Expires: session.ExpireAt,
			Path:    "/",
		})

		usr, err := user.LoadByID(ctx, tx, consumer.AuthentifiedUserID)
		if err != nil {
			return err
		}

		// Prepare http response
		resp := sdk.AuthConsumerSigninResponse{
			Token: jwt,
			User:  usr,
		}

		if err := tx.Commit(); err != nil {
			return sdk.WithStack(err)
		}

		return service.WriteJSON(w, resp, http.StatusOK)
	}
}

func (api *API) postAuthBuiltinSignoutHandler() service.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}
