package authz

import "github.com/PzKpfw-ausf-H/forgevault/internal/domain"

type Authorizer struct{}

// Both CanView / CanDownload could have asset as a parameter (as example if we had private assets) but for now it is not necessary

func (a Authorizer) CanViewAsset(actor domain.Actor) bool {
	switch actor.Role {
	case domain.RoleStudent, domain.RoleEmployee, domain.RoleAdmin, domain.RoleSuperadmin:
		return true
	default:
		return false
	}
}

func (a Authorizer) CanDownloadAsset(actor domain.Actor) bool {
	switch actor.Role {
	case domain.RoleStudent, domain.RoleEmployee, domain.RoleAdmin, domain.RoleSuperadmin:
		return true
	default:
		return false
	}
}

func (a Authorizer) CanCreateAsset(actor domain.Actor) bool {
	switch actor.Role {
	case domain.RoleEmployee, domain.RoleAdmin, domain.RoleSuperadmin:
		return true
	default:
		return false
	}
}

func (a Authorizer) CanEditAsset(actor domain.Actor, asset domain.Asset) bool {
	switch actor.Role {
	case domain.RoleAdmin, domain.RoleSuperadmin:
		return true
	case domain.RoleEmployee:
		return actor.UserID == asset.UploadedBy
	default:
		return false
	}
}

func (a Authorizer) CanDeleteAsset(actor domain.Actor, asset domain.Asset) bool {
	switch actor.Role {
	case domain.RoleAdmin, domain.RoleSuperadmin:
		return true
	case domain.RoleEmployee:
		return actor.UserID == asset.UploadedBy
	default:
		return false
	}
}

func (a Authorizer) CanAssignRole(actor domain.Actor) bool {
	return actor.Role == domain.RoleSuperadmin
}
