package service

import (
	"fmt"
	"path"
	"strings"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
)

var allowedExt = map[string]struct{}{
	//3D - related
	".fbx":   {},
	".blend": {},
	".obj":   {},
	".mtl":   {},
	".gltf":  {},
	".glb":   {},

	// textures / 2D-related
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".tga":  {},

	//audio
	".wav":  {},
	".mp3":  {},
	".ogg":  {},
	".flac": {},

	//packages
	".zip": {},
}

func validateFilename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return repo.ErrBadRequest
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return repo.ErrBadRequest
	}

	ext := strings.ToLower(path.Ext(name))
	if _, ok := allowedExt[ext]; !ok {
		return repo.ErrBadRequest
	}

	return nil
}

func buildStorageKey(assetID domain.AssetID, version int, filename string) string {
	return fmt.Sprintf("assets/%s/v%d/%s", assetID, version, filename)
}
