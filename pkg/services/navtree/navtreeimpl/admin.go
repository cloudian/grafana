package navtreeimpl

import (
	"github.com/grafana/grafana/pkg/apimachinery/identity"
	"github.com/grafana/grafana/pkg/login/social"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/accesscontrol/ssoutils"
	"github.com/grafana/grafana/pkg/services/cloudmigration"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/navtree"
	"github.com/grafana/grafana/pkg/services/serviceaccounts"
)

// nolint: gocyclo
func (s *ServiceImpl) getAdminNode(c *contextmodel.ReqContext) (*navtree.NavLink, error) {
	var configNodes []*navtree.NavLink
	ctx := c.Req.Context()
	hasAccess := ac.HasAccess(s.accessControl, c)
	authConfigUIAvailable := s.license.FeatureEnabled(social.SAMLProviderName) || s.cfg.LDAPAuthEnabled

	generalNodeLinks := []*navtree.NavLink{}
	if hasAccess(ac.OrgPreferencesAccessEvaluator) {
		generalNodeLinks = append(generalNodeLinks, &navtree.NavLink{
			Text:     "Default preferences",
			Id:       "org-settings",
			SubTitle: "Manage preferences across an organization",
			Icon:     "sliders-v-alt",
			Url:      s.cfg.AppSubURL + "/org",
		})
	}
	if hasAccess(ac.EvalPermission(ac.ActionSettingsRead, ac.ScopeSettingsAll)) {
		generalNodeLinks = append(generalNodeLinks, &navtree.NavLink{
			Text: "Settings", SubTitle: "View the settings defined in your Grafana config", Id: "server-settings", Url: s.cfg.AppSubURL + "/admin/settings", Icon: "sliders-v-alt",
		})
	}

	generalNode := &navtree.NavLink{
		Text:     "General",
		SubTitle: "Manage default preferences and settings across Grafana",
		Id:       navtree.NavIDCfgGeneral,
		Url:      s.cfg.AppSubURL + "/admin/general",
		Icon:     "shield",
		Children: generalNodeLinks,
	}

	if len(generalNode.Children) > 0 {
		configNodes = append(configNodes, generalNode)
	}

	accessNodeLinks := []*navtree.NavLink{}
	if hasAccess(ac.EvalAny(ac.EvalPermission(ac.ActionOrgUsersRead), ac.EvalPermission(ac.ActionUsersRead, ac.ScopeGlobalUsersAll))) {
		accessNodeLinks = append(accessNodeLinks, &navtree.NavLink{
			Text: "Users", SubTitle: "Manage users in Grafana", Id: "global-users", Url: s.cfg.AppSubURL + "/admin/users", Icon: "user",
		})
	}
	if hasAccess(ac.TeamsAccessEvaluator) {
		accessNodeLinks = append(accessNodeLinks, &navtree.NavLink{
			Text:     "Teams",
			Id:       "teams",
			SubTitle: "Groups of users that have common dashboard and permission needs",
			Icon:     "users-alt",
			Url:      s.cfg.AppSubURL + "/org/teams",
		})
	}
	if enableServiceAccount(s, c) {
		accessNodeLinks = append(accessNodeLinks, &navtree.NavLink{
			Text:     "Service accounts",
			Id:       "serviceaccounts",
			SubTitle: "Use service accounts to run automated workloads in Grafana",
			Icon:     "gf-service-account",
			Url:      s.cfg.AppSubURL + "/org/serviceaccounts",
		})
	}

	if s.license.FeatureEnabled("groupsync") &&
		s.features.IsEnabled(ctx, featuremgmt.FlagGroupAttributeSync) &&
		hasAccess(ac.EvalAny(
			ac.EvalPermission("groupsync.mappings:read"),
			ac.EvalPermission("groupsync.mappings:write"),
		)) {
		accessNodeLinks = append(accessNodeLinks, &navtree.NavLink{
			Text:     "External group sync",
			Id:       "groupsync",
			SubTitle: "Manage mappings of Identity Provider groups to Grafana Roles",
			Icon:     "",
			Url:      s.cfg.AppSubURL + "/admin/access/groupsync",
		})
	}

	usersNode := &navtree.NavLink{
		Text:     "Users and access",
		SubTitle: "Configure access for individual users, teams, and service accounts",
		Id:       navtree.NavIDCfgAccess,
		Url:      s.cfg.AppSubURL + "/admin/access",
		Icon:     "shield",
		Children: accessNodeLinks,
	}

	// Always append admin access as it's injected by grafana-auth-app.
	configNodes = append(configNodes, usersNode)

	if authConfigUIAvailable && hasAccess(ssoutils.EvalAuthenticationSettings(s.cfg)) ||
		hasAccess(ssoutils.OauthSettingsEvaluator(s.cfg)) {
		configNodes = append(configNodes, &navtree.NavLink{
			Text:      "Authentication",
			Id:        "authentication",
			SubTitle:  "Manage your auth settings and configure single sign-on",
			Icon:      "signin",
			IsSection: true,
			Url:       s.cfg.AppSubURL + "/admin/authentication",
		})
	}

	configNode := &navtree.NavLink{
		Id:         navtree.NavIDCfg,
		Text:       "Administration",
		SubTitle:   "Organization: " + c.GetOrgName(),
		Icon:       "cog",
		SortWeight: navtree.WeightConfig,
		Children:   configNodes,
		Url:        s.cfg.AppSubURL + "/admin",
	}

	return configNode, nil
}

func enableServiceAccount(s *ServiceImpl, c *contextmodel.ReqContext) bool {
	hasAccess := ac.HasAccess(s.accessControl, c)
	return hasAccess(serviceaccounts.AccessEvaluator)
}
