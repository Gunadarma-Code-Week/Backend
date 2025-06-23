package dto

type ResponseStatistik struct {
	TotalPeserta    string `json:"total_peserta"`
	PesertaSeminar  string `json:"peserta_seminar"`
	PesertaHackaton string `json:"peserta_hackaton"`
	PesertaCP       string `json:"peserta_cp"`
}

type Seminar struct {
	ID              int    `json:"id"`
	NamaPeserta     string `json:"nama_peserta"`
	Email           string `json:"email"`
	NomorHp         string `json:"nomor_hp"`
	Jenjang         string `json:"jenjang"`
	NamaUniversitas string `json:"nama_universitas"`
	Dokumen         string `json:"dokumen"`
	Status          bool   `json:"status"`
}

type Anggota struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	University string `json:"university"`
}

type ResponseSeminar struct {
	Seminar []Seminar `json:"seminar"`
	HasMore bool      `json:"has_more"`
}

type Hackaton struct {
	ID           int       `json:"id"`
	NamaTim      string    `json:"nama_tim"`
	JoinCode     string    `json:"join_code"`
	Members      []Anggota `json:"members"`
	KomitmenFee  string    `json:"komitmen_fee"`
	ProposalUrl  string    `json:"proposal_url"`
	PitchDeckUrl string    `json:"pitch_deck_url"`
	GithubUrl    string    `json:"github_url"`
	Stage        string    `json:"stage"`
	Status       string    `json:"status"`
}

type ResponseHackaton struct {
	Hackaton []Hackaton `json:"hackaton"`
	HasMore  bool       `json:"has_more"`
}

type Cp struct {
	ID          int       `json:"id"`
	NamaTim     string    `json:"nama_tim"`
	JoinCode    string    `json:"join_code"`
	Members     []Anggota `json:"members"`
	KomitmenFee string    `json:"komitmen_fee"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Stage       string    `json:"stage"`
	Status      string    `json:"status"`
}

type ResponseCp struct {
	Cp      []Cp `json:"cp"`
	HasMore bool `json:"has_more"`
}
