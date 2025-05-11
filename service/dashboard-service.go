package service

import (
	"gcw/dto"
	"gcw/entity"

	"gorm.io/gorm"
)

type DashboardServices struct {
	DB *gorm.DB
}

func NewDashboardServices(db *gorm.DB) DashboardServices {
	return DashboardServices{
		DB: db,
	}
}

func FindUserById(id string) {}

func (s *DashboardServices) GetAllSeminar(count, page int) ([]dto.ResponseSeminar, error) {
	var dataSeminars []entity.Seminar

	offset := page * count

	if err := s.DB.Preload("User").
		Limit(count + 1).
		Offset(offset).
		Find(&dataSeminars).Error; err != nil {
		return []dto.ResponseSeminar{}, err
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseData []dto.ResponseSeminar

	for _, data := range dataSeminars {
		dataSeminar := dto.Seminar{
			ID:              int(data.ID_Seminar),
			NamaPeserta:     data.User.Name,
			Email:           data.User.Email,
			NomorHp:         data.User.Phone,
			Jenjang:         data.User.Jenjang,
			NamaUniversitas: data.User.Institusi,
			Dokumen:         data.User.DokumenFilename,
			Status:          data.User.DataHasVerified,
		}

		response := dto.ResponseSeminar{
			Seminar: dataSeminar,
			HasMore: hasMore,
		}

		responseData = append(responseData, response)
	}

	return responseData, nil
}

func (s *DashboardServices) GetAllHackaton(count, page int) ([]dto.ResponseHackaton, error) {
	var dataSeminars []entity.HackathonTeam

	offset := page * count

	if err := s.DB.Preload("Team").
		Limit(count + 1).
		Offset(offset).
		Find(&dataSeminars).Error; err != nil {
		return []dto.ResponseHackaton{}, err
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseData []dto.ResponseHackaton

	for _, data := range dataSeminars {
		var Leader entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
			return []dto.ResponseHackaton{}, err
		}

		dataHackaton := dto.Hackaton{
			ID:           int(data.ID_HackathonTeam),
			NamaTim:      data.Team.TeamName,
			Leader:       Leader.Name,
			KomitmenFee:  data.Team.KomitmenFee,
			ProposalUrl:  data.ProposalUrl,
			PitchDeckUrl: data.PitchDeckUrl,
			GithubUrl:    data.GithubProjectUrl,
			Stage:        data.Stage,
		}

		var anggota []entity.User
		if err := s.DB.Where("id_team = ?", dataHackaton.ID).Find(&anggota).Error; err != nil {
			return []dto.ResponseHackaton{}, err
		}

		var anggotas []dto.Anggota
		for _, data := range anggota {
			if data.ID == Leader.ID {
				continue
			}

			anggota := dto.Anggota{
				Name: data.Name,
			}

			anggotas = append(anggotas, anggota)
		}

		dataHackaton.Anggota = anggotas

		response := dto.ResponseHackaton{
			Hackaton: dataHackaton,
			HasMore:  hasMore,
		}

		responseData = append(responseData, response)
	}

	return responseData, nil
}

func (s *DashboardServices) GetAllCp(count, page int) ([]dto.ResponseCp, error) {
	var dataSeminars []entity.CPTeam

	offset := page * count

	if err := s.DB.Preload("Team").
		Limit(count + 1).
		Offset(offset).
		Find(&dataSeminars).Error; err != nil {
		return []dto.ResponseCp{}, err
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseData []dto.ResponseCp

	for _, data := range dataSeminars {
		var Leader entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
			return []dto.ResponseCp{}, err
		}

		dataCp := dto.Cp{
			ID:          int(data.ID_CPTeam),
			NamaTim:     data.Team.TeamName,
			Leader:      Leader.Name,
			KomitmenFee: data.Team.KomitmenFee,
			Email:       Leader.Email,
			Password:    data.DomjudgeUsername,
			Stage:       data.Stage,
		}

		var anggota []entity.User
		if err := s.DB.Where("id_team = ?", dataCp.ID).Find(&anggota).Error; err != nil {
		}

		var anggotas []dto.Anggota
		for _, dataAnggota := range anggota {
			if dataAnggota.ID == Leader.ID {
				continue
			}

			anggota := dto.Anggota{
				Name: dataAnggota.Name,
			}

			anggotas = append(anggotas, anggota)
		}

		dataCp.Anggota = anggotas

		response := dto.ResponseCp{
			Cp:      dataCp,
			HasMore: hasMore,
		}

		responseData = append(responseData, response)
	}

	return responseData, nil
}

func (s *DashboardServices) GetEventSevice(id_user string) (dto.ResponseEvents, error) {
	var user entity.User
	if err := s.DB.Preload("Team").Where("id = ?", id_user).First(&user).Error; err != nil {
		return dto.ResponseEvents{}, err
	}

	// ambil semua member tim
	var member []entity.User
	if err := s.DB.Where("id_team = ?", user.IDTeam).Find(&member).Error; err != nil {
		return dto.ResponseEvents{}, err
	}

	// buat response user
	responseUser := dto.User{
		Name:       user.Name,
		Email:      user.Email,
		University: user.Institusi,
		Degree:     user.Jenjang,
	}

	// buat response anggota
	var responseMembers []dto.Member
	for _, m := range member {
		role := "Member"
		if m.ID == user.Team.ID_LeadTeam {
			role = "Leader"
		}
		responseMembers = append(responseMembers, dto.Member{
			Name:  m.Name,
			Role:  role,
			Email: m.Email,
		})
	}

	// ambil data event
	var hackaton entity.HackathonTeam
	var cp entity.CPTeam
	var seminar entity.Seminar

	_ = s.DB.Where("id_team = ?", user.IDTeam).First(&hackaton)
	_ = s.DB.Where("id_team = ?", user.IDTeam).First(&cp)
	_ = s.DB.Where("id_user = ?", id_user).First(&seminar)

	seminarStatus := "Unregistered"
	if seminar.ID_Seminar != 0 {
		seminarStatus = "Registered"
	}
	// cpStatus := "Unregistered"
	// if cp.ID_CPTeam != 0 {
	// 	cpStatus = "Registered"
	// }
	// hackatonStatus := "Unregistered"
	// if hackaton.ID_HackathonTeam != 0 {
	// 	hackatonStatus = "Registered"
	// }

	responseTeam := dto.Team{
		TeamName: user.Team.TeamName,
		Members:  responseMembers,
	}

	// isi data event
	events := []dto.Event{
		{
			Name:          "Seminar Nasional Teknologi AI",
			Status:        seminarStatus,
			PaymentStatus: seminar.PaymentStatus,
			Ticket:        dto.Ticket{},
		},
		{
			Name:          "Hackathon",
			Status:        hackaton.Stage,
			PaymentStatus: user.Team.KomitmenFee,
			Team:          responseTeam,
		},
		{
			Name:          "Competitive Programming",
			Status:        cp.Stage,
			PaymentStatus: user.Team.KomitmenFee,
			Team:          responseTeam,
		},
	}

	return dto.ResponseEvents{
		User:   responseUser,
		Events: events,
		IdTeam: user.Team.JoinCode,
	}, nil
}

func (s *DashboardServices) DeletePesertaService(acara, id_user string) (string, error) {
	var idUser string
	switch acara {
	case "seminar":
		var data entity.Seminar
		if err := s.DB.Where("id_user = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}

		data.IsDeleted = true

		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}

	case "hackaton":
		var data entity.HackathonTeam
		if err := s.DB.Where("id_user = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}

		data.IsDeleted = true

		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}

	case "cp":
		var data entity.CPTeam
		if err := s.DB.Where("id_user = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}

		data.IsDeleted = true

		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}

	}

	return idUser, nil
}

func (s *DashboardServices) UpdateSeminarService(id string, input dto.Seminar) error {
	var seminar entity.Seminar
	if err := s.DB.Preload("User").Where("id_seminar = ?", id).First(&seminar).Error; err != nil {
		return err
	}

	seminar.User.Name = input.NamaPeserta
	seminar.User.Email = input.Email
	seminar.User.Phone = input.NomorHp
	seminar.User.Jenjang = input.Jenjang
	seminar.User.Institusi = input.NamaUniversitas
	seminar.User.DokumenFilename = input.Dokumen
	seminar.User.DataHasVerified = input.Status

	if err := s.DB.Save(&seminar).Error; err != nil {
		return err
	}

	return nil
}

func (s *DashboardServices) UpdateHackatonService(id string, input dto.Hackaton) error {
	var hackaton entity.HackathonTeam
	if err := s.DB.Preload("Team").Where("id_hackaton_team = ?", id).First(&hackaton).Error; err != nil {
		return err
	}

	hackaton.Team.TeamName = input.NamaTim
	hackaton.ProposalUrl = input.ProposalUrl
	hackaton.GithubProjectUrl = input.GithubUrl
	hackaton.PitchDeckUrl = input.PitchDeckUrl
	hackaton.Stage = input.Stage

	if err := s.DB.Save(&hackaton).Error; err != nil {
		return err
	}

	return nil
}

func (s *DashboardServices) UpdateCpService(id string, input dto.Cp) error {
	var cp entity.CPTeam
	if err := s.DB.Preload("Team").Where("id_hackaton_team = ?", id).First(&cp).Error; err != nil {
		return err
	}

	cp.Team.TeamName = input.NamaTim
	cp.Stage = input.Stage

	if err := s.DB.Save(&cp).Error; err != nil {
		return err
	}

	return nil
}
