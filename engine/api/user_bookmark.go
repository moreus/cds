package api

import (
	"context"
	"net/http"

	"github.com/ovh/cds/engine/api/project"
	"github.com/ovh/cds/engine/api/workflow"
	"github.com/ovh/cds/engine/service"
	"github.com/ovh/cds/sdk"
)

// postUserFavoriteHandler post favorite user for workflow or project
func (api *API) postUserFavoriteHandler() service.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		params := sdk.FavoriteParams{}
		if err := service.UnmarshalBody(r, &params); err != nil {
			return err
		}

		u := getAuthentifiedUser(ctx)

		p, err := project.Load(api.mustDB(), api.Cache, params.ProjectKey, u, project.LoadOptions.WithIntegrations, project.LoadOptions.WithFavorites(u.OldUserStruct.ID))
		if err != nil {
			return sdk.WrapError(err, "unable to load projet")
		}

		switch params.Type {
		case "workflow":
			wf, errW := workflow.Load(ctx, api.mustDB(), api.Cache, p, params.WorkflowName, u, workflow.LoadOptions{})
			if errW != nil {
				return sdk.WrapError(errW, "postUserFavoriteHandler> Cannot load workflow %s/%s", params.ProjectKey, params.WorkflowName)
			}

			wf.Favorite, errW = workflow.IsFavorite(api.mustDB(), wf, u.OldUserStruct.ID)
			if errW != nil {
				return sdk.WrapError(errW, "postUserFavoriteHandler> Cannot load workflow favorite %s/%s", params.ProjectKey, params.WorkflowName)
			}
			if err := workflow.UpdateFavorite(api.mustDB(), wf.ID, u.OldUserStruct.ID, !wf.Favorite); err != nil {
				return sdk.WrapError(err, "Cannot change workflow %s/%s favorite", params.ProjectKey, params.WorkflowName)
			}
			wf.Favorite = !wf.Favorite

			return service.WriteJSON(w, wf, http.StatusOK)
		case "project":
			if err := project.UpdateFavorite(api.mustDB(), p.ID, u.OldUserStruct.ID, !p.Favorite); err != nil {
				return sdk.WrapError(err, "Cannot change workflow %s favorite", p.Key)
			}
			p.Favorite = !p.Favorite

			return service.WriteJSON(w, p, http.StatusOK)
		}

		return sdk.ErrInvalidFavoriteType
	}
}
