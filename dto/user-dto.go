package dto

type UserResponseDTO struct {
	ID                uint64 `json:"id"`
	Email             string `json:"email" binding:"max=100"`
	Role              string `json:"role" binding:"max=50"`
	Name              string `json:"name" binding:"max=50"`
	Major             string `json:"major" binding:"max=50"`
	ProfilePicture    string `json:"profile_picture" binding:"max=255"`
	NIM               string `json:"nim" binding:"max=50"`
	Institusi         string `json:"institusi" binding:"max=50"`
	Phone             string `json:"phone" binding:"max=50"`
	SocMedDocument    string `json:"socmed_document" binding:"max=255"`
	DokumenFilename   string `json:"dokumen_filename" binding:"max=255"`
	ProfileHasUpdated bool   `json:"profile_has_updated"`
	DataHasVerified   bool   `json:"data_has_verified"`
	HasPassword       bool   `json:"has_password"`
}

type UpdateUserProfileDTO struct {
	Name       string `json:"name" binding:"required,max=50"`
	NIM        string `json:"nim" binding:"required,max=50"`
	Phone      string `json:"phone" binding:"required,max=50"`
	Major      string `json:"major" binding:"required,max=50"`
	Institusi  string `json:"institusi" binding:"required,max=50"`

	SocMedDocument string `json:"socmed_document" binding:"required,max=255"`
}

// Admin User Management DTOs
type AdminGetUsersQueryDTO struct {
	Page      int    `form:"page" binding:"min=1" json:"page"`
	Limit     int    `form:"limit" binding:"min=1,max=100" json:"limit"`
	StartDate string `form:"startDate" json:"startDate" binding:"max=50"`
	EndDate   string `form:"endDate" json:"endDate" binding:"max=50"`
	Q         string `form:"q" json:"q" binding:"max=50"`
	SortBy    string `form:"sortBy" json:"sortBy" binding:"max=50"`
	SortOrder string `form:"sortOrder" binding:"oneof=ASC DESC,max=50" json:"sortOrder"`
}

type AdminUpdateUserDTO struct {
	Name              string `json:"name" binding:"max=50"`
	Email             string `json:"email" binding:"max=100"`
	Role              string `json:"role" binding:"omitempty,oneof=user admin superadmin,max=50"`
	Institusi         string `json:"institusi" binding:"max=50"`
	Phone             string `json:"phone" binding:"max=50"`
	Jenjang           string `json:"jenjang" binding:"max=50"`
	Major             string `json:"major" binding:"max=50"`
	NIM               string `json:"nim" binding:"max=50"`
	SocMedDocument    string `json:"soc_med_document" binding:"max=255"`
	DokumenFilename   string `json:"dokumen_filename" binding:"max=255"`
	ProfilePicture    string `json:"profile_picture" binding:"max=255"`
	ProfileHasUpdated *bool  `json:"profile_has_updated"`
	DataHasVerified   *bool  `json:"data_has_verified"`
}

type AdminDeleteUserDTO struct {
	Alasan string `json:"alasan" binding:"required,max=50"`
}

type AdminUserResponseDTO struct {
	ID                uint64 `json:"id"`
	Email             string `json:"email" binding:"max=100"`
	Role              string `json:"role" binding:"max=50"`
	Name              string `json:"name" binding:"max=50"`
	Institusi         string `json:"institusi" binding:"max=50"`
	Phone             string `json:"phone" binding:"max=50"`
	Jenjang           string `json:"jenjang" binding:"max=50"`
	Major             string `json:"major" binding:"max=50"`
	NIM               string `json:"nim" binding:"max=50"`
	SocMedDocument    string `json:"soc_med_document" binding:"max=255"`
	DokumenFilename   string `json:"dokumen_filename" binding:"max=255"`
	ProfilePicture    string `json:"profile_picture" binding:"max=255"`
	ProfileHasUpdated bool   `json:"profile_has_updated"`
	DataHasVerified   bool   `json:"data_has_verified"`
	IDTeam            uint64 `json:"id_team"`
	TeamName          string `json:"team_name" binding:"max=50"`
	CreatedAt         string `json:"created_at" binding:"max=50"`
	UpdatedAt         string `json:"updated_at" binding:"max=50"`
}

type AdminUsersListResponseDTO struct {
	Users []AdminUserResponseDTO `json:"users"`
	Meta  AdminUsersMetaDTO      `json:"meta"`
}

type AdminUsersMetaDTO struct {
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int64 `json:"totalPages"`
	CurrentPage int   `json:"currentPage"`
	Limit       int   `json:"limit"`
	HasMore     bool  `json:"hasMore"`
}

type UserGrowthAnalyticsDTO struct {
	StartDate string `form:"startDate" binding:"required,max=50" json:"startDate"`
	EndDate   string `form:"endDate" binding:"required,max=50" json:"endDate"`
}

type UserGrowthResponseDTO struct {
	Period    string `json:"period" binding:"max=50"`
	NewUsers  int64  `json:"newUsers"`
	TotalUsers int64 `json:"totalUsers"`
}

// type RegisterDTO struct {
// 	GoogleIdToken string `json:"google_id_token"`
// }

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
