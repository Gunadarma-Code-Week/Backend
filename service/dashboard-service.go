package service

import (
	"bytes"
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"html/template"
	"os"
	"sort"

	"github.com/sahilm/fuzzy"
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

func (s *DashboardServices) GetAllSeminar(count, page int, search string) (dto.ResponseSeminar, error) {
	var dataSeminars []entity.Seminar

	offset := page * count

	query := s.DB.Preload("User").Where("is_deleted = ? OR is_deleted IS NULL", false)

	var filteredData []entity.Seminar
	if search != "" {
		if err := query.Order("id_seminar ASC").Find(&dataSeminars).Error; err != nil {
			return dto.ResponseSeminar{}, err
		}

		matches := fuzzy.FindFrom(search, seminarSource(dataSeminars))
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})

		for _, match := range matches {
			filteredData = append(filteredData, dataSeminars[match.Index])
		}

		start := offset
		if start > len(filteredData) {
			dataSeminars = []entity.Seminar{}
		} else {
			end := start + count + 1
			if end > len(filteredData) {
				end = len(filteredData)
			}
			dataSeminars = filteredData[start:end]
		}
	} else {
		if err := query.Order("id_seminar ASC").
			Limit(count + 1).
			Offset(offset).
			Find(&dataSeminars).Error; err != nil {
			return dto.ResponseSeminar{}, err
		}
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

	var total int64
	if search != "" {
		total = int64(len(filteredData))
	} else {
		// Use a clone of the query to count without pagination
		countQuery := s.DB.Model(&entity.Seminar{}).Where("is_deleted = ? OR is_deleted IS NULL", false)
		countQuery.Count(&total)
	}
	totalPages := int((total + int64(count) - 1) / int64(count))

	response := dto.ResponseSeminar{
		Seminar:    responseData,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *DashboardServices) GetAllHackaton(count, page int, search, stage string) (dto.ResponseHackaton, error) {
	var dataSeminars []entity.HackathonTeam

	offset := page * count

	query := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false)
	if stage != "" {
		query = query.Where("stage = ?", stage)
	}

	var filteredData []entity.HackathonTeam
	if search != "" {
		// Fetch all for fuzzy matching
		if err := query.Order("id_hackathon_team ASC").Find(&dataSeminars).Error; err != nil {
			return dto.ResponseHackaton{}, err
		}

		// Perform fuzzy search
		matches := fuzzy.FindFrom(search, hackathonTeamSource(dataSeminars))
		
		// Sort by score
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})

		for _, match := range matches {
			filteredData = append(filteredData, dataSeminars[match.Index])
		}

		// Manual pagination
		start := offset
		if start > len(filteredData) {
			dataSeminars = []entity.HackathonTeam{}
		} else {
			end := start + count + 1
			if end > len(filteredData) {
				end = len(filteredData)
			}
			dataSeminars = filteredData[start:end]
		}
	} else {
		if err := query.Order("id_hackathon_team ASC").
			Limit(count + 1).
			Offset(offset).
			Find(&dataSeminars).Error; err != nil {
			return dto.ResponseHackaton{}, err
		}
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseHackatons []dto.Hackaton

	for _, data := range dataSeminars {
		if data.Team.ID_Team == 0 {
			continue // Skip if team info is missing
		}
		var Leader entity.User
		if err := s.DB.Where("id = ?", data.Team.ID_LeadTeam).First(&Leader).Error; err != nil {
			continue
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
			continue
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

	var total int64
	if search != "" {
		total = int64(len(filteredData))
	} else {
		countQuery := s.DB.Model(&entity.HackathonTeam{}).Where("is_deleted = ? OR is_deleted IS NULL", false)
		if stage != "" {
			countQuery = countQuery.Where("stage = ?", stage)
		}
		countQuery.Count(&total)
	}
	totalPages := int((total + int64(count) - 1) / int64(count))

	responseData := dto.ResponseHackaton{
		Hackaton:   responseHackatons,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}

	return responseData, nil
}

func (s *DashboardServices) GetAllCp(count, page int, search, stage string) (dto.ResponseCp, error) {
	var dataSeminars []entity.CPTeam

	offset := page * count

	query := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false)
	if stage != "" {
		query = query.Where("stage = ?", stage)
	}

	var filteredData []entity.CPTeam
	if search != "" {
		if err := query.Order("id_cp_team ASC").Find(&dataSeminars).Error; err != nil {
			return dto.ResponseCp{}, err
		}

		matches := fuzzy.FindFrom(search, cpTeamSource(dataSeminars))
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})

		for _, match := range matches {
			filteredData = append(filteredData, dataSeminars[match.Index])
		}

		start := offset
		if start > len(filteredData) {
			dataSeminars = []entity.CPTeam{}
		} else {
			end := start + count + 1
			if end > len(filteredData) {
				end = len(filteredData)
			}
			dataSeminars = filteredData[start:end]
		}
	} else {
		if err := query.Order("id_cp_team ASC").
			Limit(count + 1).
			Offset(offset).
			Find(&dataSeminars).Error; err != nil {
			return dto.ResponseCp{}, err
		}
	}

	hasMore := false

	if len(dataSeminars) > count {
		hasMore = true
		dataSeminars = dataSeminars[:count]
	}

	var responseData []dto.Cp

	for _, data := range dataSeminars {
		if data.Team.ID_Team == 0 {
			continue
		}
		var Leader entity.User
		if err := s.DB.Where("id = ?", data.Team.ID_LeadTeam).First(&Leader).Error; err != nil {
			continue
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
			continue
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

	var total int64
	if search != "" {
		total = int64(len(filteredData))
	} else {
		countQuery := s.DB.Model(&entity.CPTeam{}).Where("is_deleted = ? OR is_deleted IS NULL", false)
		if stage != "" {
			countQuery = countQuery.Where("stage = ?", stage)
		}
		countQuery.Count(&total)
	}
	totalPages := int((total + int64(count) - 1) / int64(count))

	response := dto.ResponseCp{
		Cp:         responseData,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *DashboardServices) GetAllCtf(count, page int, search, stage string) (dto.ResponseCtf, error) {
	var dataCtfTeams []entity.CTFTeam

	offset := page * count

	query := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false)
	if stage != "" {
		query = query.Where("stage = ?", stage)
	}

	var filteredData []entity.CTFTeam
	if search != "" {
		if err := query.Order("id_ctf_team ASC").Find(&dataCtfTeams).Error; err != nil {
			return dto.ResponseCtf{}, err
		}

		matches := fuzzy.FindFrom(search, ctfTeamSource(dataCtfTeams))
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})

		for _, match := range matches {
			filteredData = append(filteredData, dataCtfTeams[match.Index])
		}

		start := offset
		if start > len(filteredData) {
			dataCtfTeams = []entity.CTFTeam{}
		} else {
			end := start + count + 1
			if end > len(filteredData) {
				end = len(filteredData)
			}
			dataCtfTeams = filteredData[start:end]
		}
	} else {
		if err := query.Order("id_ctf_team ASC").
			Limit(count + 1).
			Offset(offset).
			Find(&dataCtfTeams).Error; err != nil {
			return dto.ResponseCtf{}, err
		}
	}

	hasMore := false

	if len(dataCtfTeams) > count {
		hasMore = true
		dataCtfTeams = dataCtfTeams[:count]
	}

	var responseData []dto.Ctf

	for _, data := range dataCtfTeams {
		if data.Team.ID_Team == 0 {
			continue
		}
		var Leader entity.User
		if err := s.DB.Where("id = ?", data.Team.ID_LeadTeam).First(&Leader).Error; err != nil {
			continue
		}

		dataCtf := dto.Ctf{
			ID:          int(data.ID_CTFTeam),
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
			continue
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

	var total int64
	if search != "" {
		total = int64(len(filteredData))
	} else {
		countQuery := s.DB.Model(&entity.CTFTeam{}).Where("is_deleted = ? OR is_deleted IS NULL", false)
		if stage != "" {
			countQuery = countQuery.Where("stage = ?", stage)
		}
		countQuery.Count(&total)
	}
	totalPages := int((total + int64(count) - 1) / int64(count))

	response := dto.ResponseCtf{
		Ctf:        responseData,
		HasMore:    hasMore,
		TotalPages: totalPages,
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

func (s *DashboardServices) DeletePesertaService(acara, id_user, reason string) (string, error) {
	var targetName string
	switch acara {
	case "seminar":
		var data entity.Seminar
		if err := s.DB.Preload("User").Where("id_seminar = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}
		targetName = data.User.Name
		data.IsDeleted = true
		data.DeletionReason = reason
		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}

	case "hackaton", "hackathon":
		var data entity.HackathonTeam
		if err := s.DB.Preload("Team").Where("id_hackathon_team = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}
		targetName = data.Team.TeamName
		data.IsDeleted = true
		data.DeletionReason = reason
		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}

	case "cp":
		var data entity.CPTeam
		if err := s.DB.Preload("Team").Where("id_cp_team = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}
		targetName = data.Team.TeamName
		data.IsDeleted = true
		data.DeletionReason = reason
		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}

	case "ctf":
		var data entity.CTFTeam
		if err := s.DB.Preload("Team").Where("id_ctf_team = ?", id_user).First(&data).Error; err != nil {
			return "", err
		}
		targetName = data.Team.TeamName
		data.IsDeleted = true
		data.DeletionReason = reason
		if err := s.DB.Save(&data).Error; err != nil {
			return "", err
		}
	}

	return targetName, nil
}

func (s *DashboardServices) UpdateSeminarService(id string, input dto.Seminar) (string, error) {
	var seminar entity.Seminar
	if err := s.DB.Preload("User").Where("id_seminar = ?", id).First(&seminar).Error; err != nil {
		return "", err
	}

	seminar.User.Name = input.User.Name
	seminar.User.Email = input.User.Email
	seminar.User.Phone = input.User.Phone
	seminar.User.Jenjang = input.User.Jenjang
	seminar.User.Institusi = input.User.Institusi
	seminar.PaymentStatus = input.PaymentStatus

	if err := s.DB.Save(&seminar).Error; err != nil {
		return "", err
	}

	return seminar.User.Name, nil
}

func (s *DashboardServices) UpdateHackatonService(id string, input dto.Hackaton) (string, error) {
	var hackaton entity.HackathonTeam
	if err := s.DB.Preload("Team").Where("id_hackathon_team = ?", id).First(&hackaton).Error; err != nil {
		return "", err
	}

	oldStage := hackaton.Stage
	tx := s.DB.Begin()

	if input.NamaTim != "" {
		if err := tx.Model(&entity.Team{}).
			Where("id_team = ?", hackaton.Team.ID_Team).
			Update("team_name", input.NamaTim).Error; err != nil {
			tx.Rollback()
			return "", err
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
		if input.Stage == "Stage-1" || input.Stage == "Stage-2" || input.Stage == "Final" {
			if hackaton.Status == "Registration" {
				hackaton.Status = "Verified"
			}
		}
	}

	if err := tx.Save(&hackaton).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	teamName := hackaton.Team.TeamName
	if input.NamaTim != "" {
		teamName = input.NamaTim
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
				if input.Stage == "Stage-1" {
					go s.EmailService.SendEmailHTML("Registration Confirmation: GCW 2.0 Hackathon", emails, "template/email/hackathon_stage1.html", map[string]interface{}{"TeamName": hackaton.Team.TeamName})
				} else if input.Stage != "Registered" {
					go s.EmailService.SendEmailHTML("Selection Announcement: GCW 2.0 Hackathon", emails, "template/email/hackathon_announcement.html", map[string]interface{}{"TeamName": hackaton.Team.TeamName})
				}
			}
		}
	}

	return teamName, nil
}

func (s *DashboardServices) UpdateCpService(id string, input dto.Cp) (string, error) {
	var cp entity.CPTeam
	if err := s.DB.Preload("Team").Where("id_cp_team = ?", id).First(&cp).Error; err != nil {
		return "", err
	}

	oldStage := cp.Stage
	tx := s.DB.Begin()

	if input.NamaTim != "" {
		if err := tx.Model(&entity.Team{}).
			Where("id_team = ?", cp.Team.ID_Team).
			Update("team_name", input.NamaTim).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if input.Username != "" {
		cp.DomjudgeUsername = input.Username
	}
	if input.Password != "" {
		cp.DomjudgePassword = input.Password
	}
	if input.Stage != "" {
		cp.Stage = input.Stage
		if input.Stage == "Stage-1" || input.Stage == "Stage-2" || input.Stage == "Final" {
			if cp.Status == "Registration" {
				cp.Status = "Verified"
			}
		}
	}

	if err := tx.Save(&cp).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	teamName := cp.Team.TeamName
	if input.NamaTim != "" {
		teamName = input.NamaTim
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
				if input.Stage == "Stage-1" {
					go s.EmailService.SendEmailHTML("Registration Confirmation: GCW 2.0 Competitive Programming", emails, "template/email/cp_stage1.html", map[string]interface{}{"TeamName": cp.Team.TeamName})
				} else if input.Stage != "Registered" {
					go s.EmailService.SendEmailHTML("Congratulations! You are selected - GCW 2.0 Competitive Programming", emails, "template/email/cp_announcement.html", map[string]interface{}{"TeamName": cp.Team.TeamName})
				}
			}
		}
	}

	return teamName, nil
}

func (s *DashboardServices) UpdateCtfService(id string, input dto.Ctf) (string, error) {
	var ctf entity.CTFTeam
	if err := s.DB.Preload("Team").Where("id_ctf_team = ?", id).First(&ctf).Error; err != nil {
		return "", err
	}

	oldStage := ctf.Stage
	tx := s.DB.Begin()

	if input.NamaTim != "" {
		if err := tx.Model(&entity.Team{}).
			Where("id_team = ?", ctf.Team.ID_Team).
			Update("team_name", input.NamaTim).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if input.Username != "" {
		ctf.DomjudgeUsername = input.Username
	}
	if input.Password != "" {
		ctf.DomjudgePassword = input.Password
	}
	if input.Stage != "" {
		ctf.Stage = input.Stage
		if input.Stage == "Stage-1" || input.Stage == "Stage-2" || input.Stage == "Final" {
			if ctf.Status == "Registration" {
				ctf.Status = "Verified"
			}
		}
	}

	if err := tx.Save(&ctf).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	teamName := ctf.Team.TeamName
	if input.NamaTim != "" {
		teamName = input.NamaTim
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
				if input.Stage == "Stage-1" {
					go s.EmailService.SendEmailHTML("Registration Confirmation: GCW 2.0 Capture The Flag", emails, "template/email/ctf_stage1.html", map[string]interface{}{"TeamName": ctf.Team.TeamName})
				} else if input.Stage != "Registered" {
					go s.EmailService.SendEmailHTML("Congratulations! You are selected - GCW 2.0 Capture The Flag", emails, "template/email/ctf_announcement.html", map[string]interface{}{"TeamName": ctf.Team.TeamName})
				}
			}
		}
	}

	return teamName, nil
}

// Fuzzy Search Sources

type hackathonTeamSource []entity.HackathonTeam

func (s hackathonTeamSource) String(i int) string {
	return s[i].Team.TeamName + " " + s[i].Team.JoinCode
}

func (s hackathonTeamSource) Len() int {
	return len(s)
}

type cpTeamSource []entity.CPTeam

func (s cpTeamSource) String(i int) string {
	return s[i].Team.TeamName + " " + s[i].Team.JoinCode
}

func (s cpTeamSource) Len() int {
	return len(s)
}

type ctfTeamSource []entity.CTFTeam

func (s ctfTeamSource) String(i int) string {
	return s[i].Team.TeamName + " " + s[i].Team.JoinCode
}

func (s ctfTeamSource) Len() int {
	return len(s)
}

type seminarSource []entity.Seminar

func (s seminarSource) String(i int) string {
	return s[i].ID_Tiket + " " + s[i].User.Name + " " + s[i].User.Email
}

func (s seminarSource) Len() int {
	return len(s)
}

func (s *DashboardServices) GetTeamNames(event string) ([]string, error) {
	var names []string
	switch event {
	case "hackaton", "hackathon":
		var teams []entity.HackathonTeam
		if err := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false).Find(&teams).Error; err != nil {
			return nil, err
		}
		for _, t := range teams {
			if t.Team.TeamName != "" {
				names = append(names, t.Team.TeamName)
			}
		}
	case "cp":
		var teams []entity.CPTeam
		if err := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false).Find(&teams).Error; err != nil {
			return nil, err
		}
		for _, t := range teams {
			if t.Team.TeamName != "" {
				names = append(names, t.Team.TeamName)
			}
		}
	case "ctf":
		var teams []entity.CTFTeam
		if err := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false).Find(&teams).Error; err != nil {
			return nil, err
		}
		for _, t := range teams {
			if t.Team.TeamName != "" {
				names = append(names, t.Team.TeamName)
			}
		}
	case "seminar":
		var participants []entity.Seminar
		if err := s.DB.Preload("User").Where("is_deleted = ? OR is_deleted IS NULL", false).Find(&participants).Error; err != nil {
			return nil, err
		}
		for _, p := range participants {
			if p.User.Name != "" {
				names = append(names, p.User.Name)
			}
		}
	}
	sort.Strings(names)
	return names, nil
}

func (s *DashboardServices) SendBulkEmail(event, stage string, targetRole string, subject, content string) error {
	// Parse base template
	baseTmplBytes, err := os.ReadFile("template/email/bulk_announcement.html")
	hasTemplate := false
	var tmpl *template.Template
	if err == nil {
		tmpl, err = template.New("bulk-email").Parse(string(baseTmplBytes))
		if err == nil {
			hasTemplate = true
		}
	}

	// We will collect a slice of jobs to process asynchronously
	type EmailJob struct {
		To      []string
		Content string
	}
	var jobs []EmailJob

	switch event {
	case "hackaton", "hackathon":
		var teams []entity.HackathonTeam
		query := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false)
		if stage != "" && stage != "all" && stage != "All" {
			query = query.Where("stage = ?", stage)
		}
		if err := query.Find(&teams).Error; err != nil {
			return err
		}
		
		for _, t := range teams {
			var users []entity.User
			if targetRole == "leader" {
				if err := s.DB.Where("id = ?", t.Team.ID_LeadTeam).Find(&users).Error; err != nil {
					return err
				}
			} else {
				if err := s.DB.Where("id_team = ?", t.Team.ID_Team).Find(&users).Error; err != nil {
					return err
				}
			}
			
			var teamEmails []string
			for _, u := range users {
				if u.Email != "" {
					teamEmails = append(teamEmails, u.Email)
				}
			}
			
			if len(teamEmails) == 0 {
				continue
			}

			// Render personalization
			renderedContent := content
			if hasTemplate {
				buf := new(bytes.Buffer)
				if err := tmpl.Execute(buf, map[string]interface{}{
					"TeamName": t.Team.TeamName, 
					"Name": t.Team.TeamName,
					"Subject": subject,
					"Content": content,
				}); err == nil {
					renderedContent = buf.String()
				}
			}
			
			jobs = append(jobs, EmailJob{To: teamEmails, Content: renderedContent})
		}

	case "cp":
		var teams []entity.CPTeam
		query := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false)
		if stage != "" && stage != "all" && stage != "All" {
			query = query.Where("stage = ?", stage)
		}
		if err := query.Find(&teams).Error; err != nil {
			return err
		}
		
		for _, t := range teams {
			var users []entity.User
			if targetRole == "leader" {
				if err := s.DB.Where("id = ?", t.Team.ID_LeadTeam).Find(&users).Error; err != nil {
					return err
				}
			} else {
				if err := s.DB.Where("id_team = ?", t.Team.ID_Team).Find(&users).Error; err != nil {
					return err
				}
			}
			
			var teamEmails []string
			for _, u := range users {
				if u.Email != "" {
					teamEmails = append(teamEmails, u.Email)
				}
			}
			
			if len(teamEmails) == 0 {
				continue
			}

			// Render personalization
			renderedContent := content
			if hasTemplate {
				buf := new(bytes.Buffer)
				if err := tmpl.Execute(buf, map[string]interface{}{
					"TeamName": t.Team.TeamName, 
					"Name": t.Team.TeamName,
					"Subject": subject,
					"Content": content,
				}); err == nil {
					renderedContent = buf.String()
				}
			}
			
			jobs = append(jobs, EmailJob{To: teamEmails, Content: renderedContent})
		}

	case "ctf":
		var teams []entity.CTFTeam
		query := s.DB.Preload("Team").Where("is_deleted = ? OR is_deleted IS NULL", false)
		if stage != "" && stage != "all" && stage != "All" {
			query = query.Where("stage = ?", stage)
		}
		if err := query.Find(&teams).Error; err != nil {
			return err
		}
		
		for _, t := range teams {
			var users []entity.User
			if targetRole == "leader" {
				if err := s.DB.Where("id = ?", t.Team.ID_LeadTeam).Find(&users).Error; err != nil {
					return err
				}
			} else {
				if err := s.DB.Where("id_team = ?", t.Team.ID_Team).Find(&users).Error; err != nil {
					return err
				}
			}
			
			var teamEmails []string
			for _, u := range users {
				if u.Email != "" {
					teamEmails = append(teamEmails, u.Email)
				}
			}
			
			if len(teamEmails) == 0 {
				continue
			}

			// Render personalization
			renderedContent := content
			if hasTemplate {
				buf := new(bytes.Buffer)
				if err := tmpl.Execute(buf, map[string]interface{}{
					"TeamName": t.Team.TeamName, 
					"Name": t.Team.TeamName,
					"Subject": subject,
					"Content": content,
				}); err == nil {
					renderedContent = buf.String()
				}
			}
			
			jobs = append(jobs, EmailJob{To: teamEmails, Content: renderedContent})
		}

	case "seminar":
		var participants []entity.Seminar
		query := s.DB.Preload("User").Where("is_deleted = ? OR is_deleted IS NULL", false)
		if err := query.Find(&participants).Error; err != nil {
			return err
		}
		
		for _, p := range participants {
			if p.User.Email == "" {
				continue
			}

			// Render personalization
			renderedContent := content
			if hasTemplate {
				buf := new(bytes.Buffer)
				if err := tmpl.Execute(buf, map[string]interface{}{
					"TeamName": p.User.Name, 
					"Name": p.User.Name,
					"Subject": subject,
					"Content": content,
				}); err == nil {
					renderedContent = buf.String()
				}
			}
			
			jobs = append(jobs, EmailJob{To: []string{p.User.Email}, Content: renderedContent})
		}
	}
	
	if len(jobs) == 0 {
		return fmt.Errorf("no recipients found")
	}
	
	// Process jobs in background
	go func() {
		for _, job := range jobs {
			_ = s.EmailService.SendEmailHTMLDirect(subject, job.To, job.Content)
		}
	}()
	
	return nil
}

type EmailTemplate struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (s *DashboardServices) GetEmailTemplates() ([]EmailTemplate, error) {
	templateFiles := []struct {
		Name    string
		Path    string
		Subject string
	}{
		{"hackathon_stage1.html", "template/email/hackathon_stage1.html", "Registration Confirmation: GCW 2.0 Hackathon"},
		{"hackathon_announcement.html", "template/email/hackathon_announcement.html", "Selection Announcement: GCW 2.0 Hackathon"},
		{"cp_stage1.html", "template/email/cp_stage1.html", "Registration Confirmation: GCW 2.0 Competitive Programming"},
		{"cp_announcement.html", "template/email/cp_announcement.html", "Congratulations! You are selected - GCW 2.0 Competitive Programming"},
		{"ctf_stage1.html", "template/email/ctf_stage1.html", "Registration Confirmation: GCW 2.0 Capture The Flag"},
		{"ctf_announcement.html", "template/email/ctf_announcement.html", "Congratulations! You are selected - GCW 2.0 Capture The Flag"},
	}

	var templates []EmailTemplate
	for _, f := range templateFiles {
		contentBytes, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		templates = append(templates, EmailTemplate{
			Name:    f.Name,
			Subject: f.Subject,
			Content: string(contentBytes),
		})
	}

	return templates, nil
}
