package dto

type ResponseStatistik struct {
	TotalPeserta    string `json:"total_peserta" binding:"max=50"`
	PesertaSeminar  string `json:"peserta_seminar" binding:"max=50"`
	PesertaHackaton string `json:"peserta_hackaton" binding:"max=50"`
	PesertaCP       string `json:"peserta_cp" binding:"max=50"`
}

type Seminar struct {
	IDSeminar     uint64    `json:"id_seminar"`
	IDTiket       string    `json:"id_tiket" binding:"max=50"`
	PaymentStatus string    `json:"payment_status" binding:"max=50"`
	User          UserInfo  `json:"user"`
	CreatedAt     string    `json:"created_at" binding:"max=50"`
	UpdatedAt     string    `json:"updated_at" binding:"max=50"`
}

// UserInfo sudah dideklarasi di seminar-dto.go

type Anggota struct {
	Name       string `json:"name" binding:"max=50"`
	Email      string `json:"email" binding:"max=100"`
	Role       string `json:"role" binding:"max=50"`
	University string `json:"university" binding:"max=50"`
}

type ResponseSeminar struct {
	Seminar    []Seminar `json:"seminar"`
	HasMore    bool      `json:"has_more"`
	TotalPages int       `json:"total_pages"`
}

type Hackaton struct {
	ID           int       `json:"id"`
	NamaTim      string    `json:"nama_tim" binding:"max=50"`
	JoinCode     string    `json:"join_code" binding:"max=50"`
	Members      []Anggota `json:"members"`
	KomitmenFee  string    `json:"komitmen_fee" binding:"max=255"`
	ProposalUrl  string    `json:"proposal_url"`
	PitchDeckUrl string    `json:"pitch_deck_url"`
	GithubUrl    string    `json:"github_url"`
	Stage        string    `json:"stage" binding:"max=50"`
	Status       string    `json:"status" binding:"max=50"`
}

type ResponseHackaton struct {
	Hackaton   []Hackaton `json:"hackaton"`
	HasMore    bool       `json:"has_more"`
	TotalPages int        `json:"total_pages"`
}

type Cp struct {
	ID          int       `json:"id"`
	NamaTim     string    `json:"nama_tim" binding:"max=50"`
	JoinCode    string    `json:"join_code" binding:"max=50"`
	Members     []Anggota `json:"members"`
	KomitmenFee string    `json:"komitmen_fee" binding:"max=255"`
	Username    string    `json:"username" binding:"max=50"`
	Password    string    `json:"password" binding:"max=50"`
	Stage       string    `json:"stage" binding:"max=50"`
	Status      string    `json:"status" binding:"max=50"`
}

type ResponseCp struct {
	Cp         []Cp `json:"cp"`
	HasMore    bool `json:"has_more"`
	TotalPages int  `json:"total_pages"`
}

type Ctf struct {
	ID          int       `json:"id"`
	NamaTim     string    `json:"nama_tim" binding:"max=50"`
	JoinCode     string    `json:"join_code" binding:"max=50"`
	Members     []Anggota `json:"members"`
	KomitmenFee string    `json:"komitmen_fee" binding:"max=255"`
	Username    string    `json:"username" binding:"max=50"`
	Password    string    `json:"password" binding:"max=50"`
	Stage       string    `json:"stage" binding:"max=50"`
	Status      string    `json:"status" binding:"max=50"`
}

type ResponseCtf struct {
	Ctf        []Ctf `json:"ctf"`
	HasMore    bool  `json:"has_more"`
	TotalPages int   `json:"total_pages"`
}

type DeleteTeamRequest struct {
	Alasan string `json:"alasan" binding:"required,max=50"`
}
