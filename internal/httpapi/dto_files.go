package httpapi

type UploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

type UploadURLResponse struct {
	Version    int    `json:"version"`
	StorageKey string `json:"storageKey"`
	UploadURL  string `json:"uploadUrl"`
	ExpiresIn  int64  `json:"expiresIn"`
}

type ConfirmRequest struct {
	Version     int    `json:"version"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	SizeBytes   int64  `json:"sizeBytes"`
	StorageKey  string `json:"storageKey"`
	Checksum    string `json:"checksum"`
}

type DownloadURLResponse struct {
	DownloadURL string `json:"downloadUrl"`
	ExpiresIn   int64  `json:"expiresIn"`
}

func toUploadResponse(version int, storageKey, uploadURL string, expiresInSec int64) UploadURLResponse {
	var ur UploadURLResponse
	ur.Version = version
	ur.StorageKey = storageKey
	ur.UploadURL = uploadURL
	ur.ExpiresIn = expiresInSec

	return ur
}

func toDownloadResponse(downloadUrl string, expiresInSec int64) DownloadURLResponse {
	var dr DownloadURLResponse
	dr.DownloadURL = downloadUrl
	dr.ExpiresIn = expiresInSec

	return dr
}
