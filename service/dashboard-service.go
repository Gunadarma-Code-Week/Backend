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
		if err := s.DB.Where("IDTeam = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
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
		if err := s.DB.Where("IDTeam = ?", dataHackaton.ID).Find(&anggota).Error; err != nil {
			return []dto.ResponseHackaton{}, err
		}

		if len(anggota) > 0 {
			dataHackaton.Anggota1 = anggota[0].Name
		}
		if len(anggota) > 1 {
			dataHackaton.Anggota2 = anggota[1].Name
		}
		if len(anggota) > 2 {
			dataHackaton.Anggota3 = anggota[2].Name
		}
		if len(anggota) > 3 {
			dataHackaton.Anggota4 = anggota[3].Name
		}

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
		if err := s.DB.Where("IDTeam = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
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
		if err := s.DB.Where("IDTeam = ?", dataCp.ID).Find(&anggota).Error; err != nil {
			return []dto.ResponseCp{}, err
		}

		if len(anggota) > 0 {
			dataCp.Anggota1 = anggota[0].Name
		}
		if len(anggota) > 1 {
			dataCp.Anggota2 = anggota[1].Name
		}

		response := dto.ResponseCp{
			Cp:      dataCp,
			HasMore: hasMore,
		}

		responseData = append(responseData, response)
	}

	return responseData, nil
}
