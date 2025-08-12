package lib

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	cldCustom "github.com/tph-kds/chat_realtime/backend/internal/database"
	// cldCustom "github.com/tph-kds/chat_realtime/backend/internal/api/handlers"
)

func UploadToCloudinary(base64Image string) (*uploader.UploadResult, error) {
	ctx := context.Background()
	cldClient := cldCustom.GetCloudinary()
	return cldClient.Upload.Upload(ctx, base64Image, uploader.UploadParams{})
}
