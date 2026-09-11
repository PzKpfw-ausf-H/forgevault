package postgres

import (
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/google/uuid"
)

func validAsset() domain.Asset {
	now := time.Now().UTC().Truncate(time.Microsecond)

	return domain.Asset{
		ID:          domain.AssetID(uuid.NewString()),
		Title:       "TestAsset",
		Description: "something to test",

		AuthorName: "Student A",
		UploadedBy: domain.UserID(uuid.NewString()),

		Type:             domain.AssetType3D,
		Tags:             []string{"tag1", "tag2"},
		ProcessingStatus: domain.ProcessingStatusReady,

		CreatedAt: now,
		UpdatedAt: now,
	}
}
