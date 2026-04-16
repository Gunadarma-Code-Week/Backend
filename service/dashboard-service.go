package service

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"os"
	"time"

	"gorm.io/gorm"
)

type DashboardServices struct {
	DB           *gorm.DB
	EmailService *EmailService
}

func NewDashboardServices(db *gorm.DB) DashboardServices {
	return DashboardServices{
		DB:           db,
		EmailService: NewEmailService(),
	}
}

func FindUserById(id string) {}

func (s *DashboardServices) GetAllSeminar(startDate, endDate time.Time, count, page int, search string) (dto.ResponseSeminar, error) {
	var dataSeminars []entity.Seminar

	offset := page * count

	query := s.DB.Preload("User").
		Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Add search functionality
	if search != "" {
		query = query.Where(
			"id_tiket LIKE ? OR EXISTS (SELECT 1 FROM users WHERE users.id = seminars.id_user AND (users.name LIKE ? OR users.email LIKE ?))",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	if err := query.
		Order("id_seminar ASC").
		Limit(count + 1).
		Offset(offset).
		Find(&dataSeminars).Error; err != nil {
		return dto.ResponseSeminar{}, err
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseData []dto.Seminar

	for _, data := range dataSeminars {
		dataSeminar := dto.Seminar{
			IDSeminar:     data.ID_Seminar,
			IDTiket:       data.ID_Tiket,
			PaymentStatus: data.PaymentStatus,
			User: dto.UserInfo{
				ID:        data.User.ID,
				Name:      data.User.Name,
				Email:     data.User.Email,
				Phone:     data.User.Phone,
				Jenjang:   data.User.Jenjang,
				Institusi: data.User.Institusi,
			},
			CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		responseData = append(responseData, dataSeminar)
	}

	response := dto.ResponseSeminar{
		Seminar: responseData,
		HasMore: hasMore,
	}

	return response, nil
}

func (s *DashboardServices) GetAllHackaton(startDate, endDate time.Time, count, page int) (dto.ResponseHackaton, error) {
	var dataSeminars []entity.HackathonTeam

	offset := page * count

	if err := s.DB.Preload("Team").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("id_hackathon_team ASC").
		Limit(count + 1).
		Offset(offset).
		Find(&dataSeminars).Error; err != nil {
		return dto.ResponseHackaton{}, err
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseHackatons []dto.Hackaton

	for _, data := range dataSeminars {
		var Leader entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
			return dto.ResponseHackaton{}, err
		}

		dataHackaton := dto.Hackaton{
			ID:           int(data.ID_HackathonTeam),
			NamaTim:      data.Team.TeamName,
			JoinCode:     data.Team.JoinCode,
			KomitmenFee:  data.Team.KomitmenFee,
			ProposalUrl:  data.ProposalUrl,
			PitchDeckUrl: data.PitchDeckUrl,
			GithubUrl:    data.GithubProjectUrl,
			Stage:        data.Stage,
			Status:       data.Status,
		}

		var anggota []entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).Find(&anggota).Error; err != nil {
			return dto.ResponseHackaton{}, err
		}

		// Create members list including leader
		var members []dto.Anggota

		// Add leader as the first member
		leaderMember := dto.Anggota{
			Name:       Leader.Name,
			Email:      Leader.Email,
			Role:       "leader",
			University: Leader.Institusi,
		}
		members = append(members, leaderMember)

		// Add other team members
		for _, data := range anggota {
			if data.ID == Leader.ID {
				continue // Skip leader as they're already added
			}

			anggota := dto.Anggota{
				Name:       data.Name,
				Email:      data.Email,
				Role:       "member",
				University: data.Institusi,
			}

			members = append(members, anggota)
		}

		dataHackaton.Members = members

		responseHackatons = append(responseHackatons, dataHackaton)
	}

	responseData := dto.ResponseHackaton{
		Hackaton: responseHackatons,
		HasMore:  hasMore,
	}

	return responseData, nil
}

func (s *DashboardServices) GetAllCp(startDate, endDate time.Time, count, page int) (dto.ResponseCp, error) {
	var dataSeminars []entity.CPTeam

	offset := page * count

	if err := s.DB.Preload("Team").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("id_cp_team ASC").
		Limit(count + 1).
		Offset(offset).
		Find(&dataSeminars).Error; err != nil {
		return dto.ResponseCp{}, err
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseData []dto.Cp

	for _, data := range dataSeminars {
		var Leader entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
			return dto.ResponseCp{}, err
		}

		dataCp := dto.Cp{
			ID:          int(data.ID_CPTeam),
			JoinCode:    data.Team.JoinCode,
			NamaTim:     data.Team.TeamName,
			KomitmenFee: data.Team.KomitmenFee,
			Username:    data.DomjudgeUsername,
			Password:    data.DomjudgePassword,
			Stage:       data.Stage,
			Status:      data.Status,
		}

		var anggota []entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).Find(&anggota).Error; err != nil {
			return dto.ResponseCp{}, err
		}

		// Create members list including leader
		var members []dto.Anggota

		// Add leader as the first member
		leaderMember := dto.Anggota{
			Name:       Leader.Name,
			Email:      Leader.Email,
			Role:       "leader",
			University: Leader.Institusi,
		}
		members = append(members, leaderMember)

		// Add other team members
		for _, dataAnggota := range anggota {
			if dataAnggota.ID == Leader.ID {
				continue // Skip leader as they're already added
			}

			anggota := dto.Anggota{
				Name:       dataAnggota.Name,
				Email:      dataAnggota.Email,
				Role:       "member",
				University: dataAnggota.Institusi,
			}

			members = append(members, anggota)
		}

		dataCp.Members = members

		responseData = append(responseData, dataCp)
	}

	response := dto.ResponseCp{
		Cp:      responseData,
		HasMore: hasMore,
	}

	return response, nil
}

func (s *DashboardServices) GetAllCtf(startDate, endDate time.Time, count, page int) (dto.ResponseCtf, error) {
	var dataCtfTeams []entity.CTFTeam

	offset := page * count

	if err := s.DB.Preload("Team").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("id_ctf_team ASC").
		Limit(count + 1).
		Offset(offset).
		Find(&dataCtfTeams).Error; err != nil {
		return dto.ResponseCtf{}, err
	}

	hasMore := false

	if len(dataCtfTeams) > count {
		hasMore = true
		dataCtfTeams = dataCtfTeams[:count]
	}

	var responseData []dto.Ctf

	for _, data := range dataCtfTeams {
		var Leader entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).First(&Leader).Error; err != nil {
			return dto.ResponseCtf{}, err
		}

		dataCtf := dto.Ctf{
			ID:          int(data.ID_CTFTeam),
			JoinCode:    data.Team.JoinCode,
			NamaTim:     data.Team.TeamName,
			KomitmenFee: data.Team.KomitmenFee,
			Stage:       data.Stage,
			Status:      data.Status,
		}

		var anggota []entity.User
		if err := s.DB.Where("id_team = ?", data.Team.ID_Team).Find(&anggota).Error; err != nil {
			return dto.ResponseCtf{}, err
		}

		var members []dto.Anggota

		leaderMember := dto.Anggota{
			Name:       Leader.Name,
			Email:      Leader.Email,
			Role:       "leader",
			University: Leader.Institusi,
		}
		members = append(members, leaderMember)

		for _, dataAnggota := range anggota {
			if dataAnggota.ID == Leader.ID {
				continue
			}

			anggota := dto.Anggota{
				Name:       dataAnggota.Name,
				Email:      dataAnggota.Email,
				Role:       "member",
				University: dataAnggota.Institusi,
			}

			members = append(members, anggota)
		}

		dataCtf.Members = members

		responseData = append(responseData, dataCtf)
	}

	response := dto.ResponseCtf{
		Ctf:     responseData,
		HasMore: hasMore,
	}

	return response, nil
}

// func (s *DashboardServices) GetEventSevice(id_user string) (dto.ResponseEvents, error) {
// 	var user entity.User
// 	if err := s.DB.Preload("Team").Where("id = ?", id_user).First(&user).Error; err != nil {
// 		return dto.ResponseEvents{}, err
// 	}

// 	// ambil semua member tim
// 	var member []entity.User
// 	if err := s.DB.Where("id_team = ?", user.IDTeam).Find(&member).Error; err != nil {
// 		return dto.ResponseEvents{}, err
// 	}

// 	// buat response user
// 	responseUser := dto.User{
// 		ID:         user.ID,
// 		Name:       user.Name,
// 		Email:      user.Email,
// 		University: user.Institusi,
// 	}

// 	// buat response anggota
// 	var responseMembers []dto.Member
// 	for _, m := range member {
// 		role := "Member"
// 		if m.ID == user.Team.ID_LeadTeam {
// 			role = "Leader"
// 		}
// 		responseMembers = append(responseMembers, dto.Member{
// 			Name:  m.Name,
// 			Role:  role,
// 			Email: m.Email,
// 		})
// 	}

// 	// ambil data event
// 	var hackaton entity.HackathonTeam
// 	var cp entity.CPTeam
// 	var seminar entity.Seminar

// 	_ = s.DB.Where("id_team = ?", user.IDTeam).First(&hackaton)
// 	_ = s.DB.Where("id_team = ?", user.IDTeam).First(&cp)
// 	_ = s.DB.Where("id_user = ?", id_user).First(&seminar)

// 	seminarStatus := "Unregistered"
// 	if seminar.ID_Seminar != 0 {
// 		seminarStatus = "Registered"
// 	}
// 	// cpStatus := "Unregistered"
// 	// if cp.ID_CPTeam != 0 {
// 	// 	cpStatus = "Registered"
// 	// }
// 	// hackatonStatus := "Unregistered"
// 	// if hackaton.ID_HackathonTeam != 0 {
// 	// 	hackatonStatus = "Registered"
// 	// }

// 	// isi data event
// 	events := []dto.Event{
// 		{
// 			Name:   "Seminar Nasional Teknologi AI",
// 			Status: seminarStatus,
// 			Ticket: dto.Ticket{},
// 		},
// 		{
// 			Name:   "Hackathon",
// 			Status: hackaton.Stage,
// 		},
// 		{
// 			Name:   "Competitive Programming",
// 			Status: cp.Stage,
// 		},
// 	}

// 	return dto.ResponseEvents{
// 		User:   responseUser,
// 		Events: events,
// 		IdTeam: user.Team.JoinCode,
// 	}, nil
// }

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

	case "ctf":
		var data entity.CTFTeam
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

	seminar.User.Name = input.User.Name
	seminar.User.Email = input.User.Email
	seminar.User.Phone = input.User.Phone
	seminar.User.Jenjang = input.User.Jenjang
	seminar.User.Institusi = input.User.Institusi
	seminar.PaymentStatus = input.PaymentStatus

	if err := s.DB.Save(&seminar).Error; err != nil {
		return err
	}

	return nil
}

func (s *DashboardServices) UpdateHackatonService(id string, input dto.Hackaton) error {
	var hackaton entity.HackathonTeam
	if err := s.DB.Preload("Team").Where("id_hackathon_team = ?", id).First(&hackaton).Error; err != nil {
		return err
	}

	oldStage := hackaton.Stage
	tx := s.DB.Begin()
	
	if input.NamaTim != "" {
		if err := tx.Model(&entity.Team{}).
			Where("id_team = ?", hackaton.Team.ID_Team).
			Update("team_name", input.NamaTim).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.ProposalUrl != "" {
		hackaton.ProposalUrl = input.ProposalUrl
	}
	if input.GithubUrl != "" {
		hackaton.GithubProjectUrl = input.GithubUrl
	}
	if input.PitchDeckUrl != "" {
		hackaton.PitchDeckUrl = input.PitchDeckUrl
	}
	if input.Stage != "" {
		hackaton.Stage = input.Stage
	}

	if err := tx.Save(&hackaton).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if oldStage != input.Stage {
		var teamMembers []entity.User
		if err := s.DB.Where("id_team = ?", hackaton.Team.ID_Team).Find(&teamMembers).Error; err == nil {
			var emails []string
			for _, m := range teamMembers {
				if m.Email != "" {
					emails = append(emails, m.Email)
				}
			}
			if len(emails) > 0 && os.Getenv("AUTO_EMAIL") == "true" {
				msg := fmt.Sprintf("Selamat Tim %s, Anda berhasil lolos ke tahap %s!", hackaton.Team.TeamName, input.Stage)
				_ = s.EmailService.SendEmail("Pembaruan Status Tim - GCW 2026", emails, msg)
			}
		}
	}

	return nil
}

func (s *DashboardServices) UpdateCpService(id string, input dto.Cp) error {
	var cp entity.CPTeam
	if err := s.DB.Preload("Team").Where("id_cp_team = ?", id).First(&cp).Error; err != nil {
		return err
	}

	oldStage := cp.Stage
	tx := s.DB.Begin()
	
	if input.NamaTim != "" {
		if err := tx.Model(&entity.Team{}).
			Where("id_team = ?", cp.Team.ID_Team).
			Update("team_name", input.NamaTim).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Stage != "" {
		cp.Stage = input.Stage
	}

	if err := tx.Save(&cp).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if oldStage != input.Stage {
		var teamMembers []entity.User
		if err := s.DB.Where("id_team = ?", cp.Team.ID_Team).Find(&teamMembers).Error; err == nil {
			var emails []string
			for _, m := range teamMembers {
				if m.Email != "" {
					emails = append(emails, m.Email)
				}
			}
			if len(emails) > 0 && os.Getenv("AUTO_EMAIL") == "true" {
				msg := fmt.Sprintf("Selamat Tim %s, Anda berhasil lolos ke tahap %s!", cp.Team.TeamName, input.Stage)
				_ = s.EmailService.SendEmail("Pembaruan Status Tim - GCW 2026", emails, msg)
			}
		}
	}

	return nil
}

func (s *DashboardServices) UpdateCtfService(id string, input dto.Ctf) error {
	var ctf entity.CTFTeam
	if err := s.DB.Preload("Team").Where("id_ctf_team = ?", id).First(&ctf).Error; err != nil {
		return err
	}

	oldStage := ctf.Stage
	tx := s.DB.Begin()
	
	if input.NamaTim != "" {
		if err := tx.Model(&entity.Team{}).
			Where("id_team = ?", ctf.Team.ID_Team).
			Update("team_name", input.NamaTim).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Stage != "" {
		ctf.Stage = input.Stage
	}

	if err := tx.Save(&ctf).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if oldStage != input.Stage {
		var teamMembers []entity.User
		if err := s.DB.Where("id_team = ?", ctf.Team.ID_Team).Find(&teamMembers).Error; err == nil {
			var emails []string
			for _, m := range teamMembers {
				if m.Email != "" {
					emails = append(emails, m.Email)
				}
			}
			if len(emails) > 0 && os.Getenv("AUTO_EMAIL") == "true" {
				msg := fmt.Sprintf("Selamat Tim %s, Anda berhasil lolos ke tahap %s!", ctf.Team.TeamName, input.Stage)
				_ = s.EmailService.SendEmail("Pembaruan Status Tim - GCW 2026", emails, msg)
			}
		}
	}

	return nil
}
