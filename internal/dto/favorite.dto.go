package dto

type FavoriteRequest struct {
	FavoriteUserID string `json:"favorite_user_id" binding:"required,uuid"`
}
