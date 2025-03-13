package dto

type RegistrationResponseDTO struct {
	Name           string `json:"name"`
	Supervisor     string `json:"supervisor"`
	SupervisorNIDN string `json:"supervisor_nidn"`
	ID_LeadTeam    string `json:"id_lead"`
	ID_Event       string `json:"id_event"`
}
