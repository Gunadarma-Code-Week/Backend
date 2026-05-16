package dto

type UserResponseDTO struct {
	ID                uint64 `json:"id"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	Name              string `json:"name"`
	Major             string `json:"major"`
	ProfilePicture    string `json:"profile_picture"`
	NIM               string `json:"nim"`
	Institusi         string `json:"institusi"`
	Phone             string `json:"phone"`
	SocMedDocument    string `json:"socmed_document"`
	DokumenFilename   string `json:"dokumen_filename"`
	ProfileHasUpdated bool   `json:"profile_has_updated"`
	DataHasVerified   bool   `json:"data_has_verified"`
	HasPassword       bool   `json:"has_password"`
}

type UpdateUserProfileDTO struct {
	Name       string `json:"name" binding:"required"`
	NIM        string `json:"nim" binding:"required"`
	Phone      string `json:"phone"`
	Major      string `json:"major"`
	Institusi  string `json:"institusi" binding:"required"`

	SocMedDocument string `json:"socmed_document"`
}

// Admin User Management DTOs
type AdminGetUsersQueryDTO struct {
	Page      int    `form:"page" binding:"min=1" json:"page"`
	Limit     int    `form:"limit" binding:"min=1,max=100" json:"limit"`
	StartDate string `form:"startDate" json:"startDate"`
	EndDate   string `form:"endDate" json:"endDate"`
	Q         string `form:"q" json:"q"`
	SortBy    string `form:"sortBy" json:"sortBy"`
	SortOrder string `form:"sortOrder" binding:"oneof=ASC DESC" json:"sortOrder"`
}

type AdminUpdateUserDTO struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Role              string `json:"role" binding:"omitempty,oneof=user admin superadmin"`
	Institusi         string `json:"institusi"`
	Phone             string `json:"phone"`
	Jenjang           string `json:"jenjang"`
	Major             string `json:"major"`
	NIM               string `json:"nim"`
	SocMedDocument    string `json:"soc_med_document"`
	DokumenFilename   string `json:"dokumen_filename"`
	ProfilePicture    string `json:"profile_picture"`
	ProfileHasUpdated *bool  `json:"profile_has_updated"`
	DataHasVerified   *bool  `json:"data_has_verified"`
}

type AdminDeleteUserDTO struct {
	Alasan string `json:"alasan" binding:"required"`
}

type AdminUserResponseDTO struct {
	ID                uint64 `json:"id"`
	Email             string `json:"email"`
	Role              string `json:"role"`
	Name              string `json:"name"`
	Institusi         string `json:"institusi"`
	Phone             string `json:"phone"`
	Jenjang           string `json:"jenjang"`
	Major             string `json:"major"`
	NIM               string `json:"nim"`
	SocMedDocument    string `json:"soc_med_document"`
	DokumenFilename   string `json:"dokumen_filename"`
	ProfilePicture    string `json:"profile_picture"`
	ProfileHasUpdated bool   `json:"profile_has_updated"`
	DataHasVerified   bool   `json:"data_has_verified"`
	IDTeam            uint64 `json:"id_team"`
	TeamName          string `json:"team_name"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
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
	StartDate string `form:"startDate" binding:"required" json:"startDate"`
	EndDate   string `form:"endDate" binding:"required" json:"endDate"`
}

type UserGrowthResponseDTO struct {
	Period    string `json:"period"`
	NewUsers  int64  `json:"newUsers"`
	TotalUsers int64 `json:"totalUsers"`
}

// type RegisterDTO struct {
// 	GoogleIdToken string `json:"google_id_token"`
// }

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
