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

type ResponseSeminar struct {
	Seminar Seminar `json:"seminar"`
	HasMore bool    `json:"has_more"`
}

type Hackaton struct {
	ID           int `json:"id"`
	NamaTim      string `json:"nama_tim"`
	Leader       string `json:"leader"`
	Anggota1     string `json:"anggota_1"`
	Anggota2     string `json:"anggota_2"`
	Anggota3     string `json:"anggota_3"`
	Anggota4     string `json:"anggota_4"`
	KomitmenFee  string `json:"komitmen_fee"`
	ProposalUrl  string `json:"proposal_url"`
	PitchDeckUrl string `json:"pitch_deck_url"`
	GithubUrl    string `json:"github_url"`
	Tahap1       bool   `json:"tahap_1"`
	Tahap2       bool   `json:"tahap_2"`
	Final        bool   `json:"final"`
}

type ResponseHackaton struct {
	Hackaton Hackaton `json:"hackaton"`
	HasMore  bool     `json:"has_more"`
}

type Cp struct {
	ID          int `json:"id"`
	NamaTim     string `json:"nama_tim"`
	Leader      string `json:"leader"`
	Anggota1    string `json:"anggota_1"`
	Anggota2    string `json:"anggota_2"`
	KomitmenFee string `json:"komitmen_fee"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Tahap1      bool   `json:"tahap_1"`
	Final       bool   `json:"final"`
}

type ResponseCp struct {
	Cp      Cp   `json:"cp"`
	HasMore bool `json:"has_more"`
}
