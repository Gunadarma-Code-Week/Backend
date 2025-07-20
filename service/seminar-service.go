package service

import (
	"errors"
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"time"

	"gorm.io/gorm"
)

type SeminarService struct {
	DB *gorm.DB
}

func NewSeminarService(db *gorm.DB) *SeminarService {
	return &SeminarService{
		DB: db,
	}
}

// JoinSeminar - User bergabung ke seminar dengan ID tiket
func (s *SeminarService) JoinSeminar(userID uint64, request dto.JoinSeminarRequest) (*dto.JoinSeminarResponse, error) {
	// Cek apakah user sudah terdaftar di seminar
	var existingSeminar entity.Seminar
	err := s.DB.Where("id_user = ? AND is_deleted = ?", userID, false).First(&existingSeminar).Error
	if err == nil {
		return nil, errors.New("user sudah terdaftar di seminar")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Generate ID tiket unik
	var idTiket string
	for {
		// Generate ID tiket dengan format: SEM + timestamp + random string
		timestamp := time.Now().Format("20060102")
		randomStr := helper.RandomString(6)
		idTiket = fmt.Sprintf("SEM%s%s", timestamp, randomStr)

		// Cek apakah ID tiket sudah ada
		var ticketCheck entity.Seminar
		err = s.DB.Where("id_tiket = ? AND is_deleted = ?", idTiket, false).First(&ticketCheck).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ID tiket belum ada, bisa digunakan
			break
		}
		if err != nil {
			return nil, err
		}
		// Jika ID tiket sudah ada, loop lagi untuk generate yang baru
	}

	// Buat record seminar baru
	seminar := entity.Seminar{
		ID_Tiket:      idTiket,
		IDUser:        userID,
		PaymentStatus: "success",
		IsDeleted:     false,
	}

	if err := s.DB.Create(&seminar).Error; err != nil {
		return nil, err
	}

	return &dto.JoinSeminarResponse{
		Message:   "Berhasil bergabung ke seminar",
		Status:    "success",
		IDTiket:   seminar.ID_Tiket,
		SeminarID: seminar.ID_Seminar,
	}, nil
}

// GetTicketDetail - Mendapatkan detail tiket seminar user
func (s *SeminarService) GetTicketDetail(userID uint64) (*dto.SeminarTicketDetail, error) {
	var seminar entity.Seminar
	err := s.DB.Preload("User").Where("id_user = ? AND is_deleted = ?", userID, false).First(&seminar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tiket seminar tidak ditemukan")
		}
		return nil, err
	}

	// Map ke response DTO
	response := &dto.SeminarTicketDetail{
		IDSeminar:     seminar.ID_Seminar,
		IDTiket:       seminar.ID_Tiket,
		PaymentStatus: seminar.PaymentStatus,
		User: dto.UserInfo{
			ID:        seminar.User.ID,
			Name:      seminar.User.Name,
			Email:     seminar.User.Email,
			Phone:     seminar.User.Phone,
			Jenjang:   seminar.User.Jenjang,
			Institusi: seminar.User.Institusi,
		},
		CreatedAt: seminar.CreatedAt,
		UpdatedAt: seminar.UpdatedAt,
	}

	return response, nil
}

// GetTicketByID - Mendapatkan detail tiket seminar berdasarkan ID tiket (untuk admin)
func (s *SeminarService) GetTicketByID(ticketID string) (*dto.SeminarTicketDetail, error) {
	var seminar entity.Seminar
	err := s.DB.Preload("User").Where("id_tiket = ? AND is_deleted = ?", ticketID, false).First(&seminar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tiket seminar tidak ditemukan")
		}
		return nil, err
	}

	// Map ke response DTO
	response := &dto.SeminarTicketDetail{
		IDSeminar:     seminar.ID_Seminar,
		IDTiket:       seminar.ID_Tiket,
		PaymentStatus: seminar.PaymentStatus,
		User: dto.UserInfo{
			ID:        seminar.User.ID,
			Name:      seminar.User.Name,
			Email:     seminar.User.Email,
			Phone:     seminar.User.Phone,
			Jenjang:   seminar.User.Jenjang,
			Institusi: seminar.User.Institusi,
		},
		CreatedAt: seminar.CreatedAt,
		UpdatedAt: seminar.UpdatedAt,
	}

	return response, nil
}

// AdminAddParticipant - Admin menambahkan participant ke seminar berdasarkan user ID
func (s *SeminarService) AdminAddParticipant(userID uint64) (*dto.AdminAddParticipantResponse, error) {
	// Cek apakah user ada dan data sudah terverifikasi
	var user entity.User
	err := s.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}

	// Cek apakah data user sudah terverifikasi
	if !user.DataHasVerified {
		return nil, errors.New("data user belum terverifikasi")
	}

	// Cek apakah user sudah terdaftar di seminar
	var existingSeminar entity.Seminar
	err = s.DB.Where("id_user = ? AND is_deleted = ?", userID, false).First(&existingSeminar).Error
	if err == nil {
		return nil, errors.New("user sudah terdaftar di seminar")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Generate ID tiket unik
	var idTiket string
	for {
		// Generate ID tiket dengan format: SEM + timestamp + random string
		timestamp := time.Now().Format("20060102")
		randomStr := helper.RandomString(6)
		idTiket = fmt.Sprintf("SEM%s%s", timestamp, randomStr)

		// Cek apakah ID tiket sudah ada
		var ticketCheck entity.Seminar
		err = s.DB.Where("id_tiket = ? AND is_deleted = ?", idTiket, false).First(&ticketCheck).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ID tiket belum ada, bisa digunakan
			break
		}
		if err != nil {
			return nil, err
		}
		// Jika ID tiket sudah ada, loop lagi untuk generate yang baru
	}

	// Buat record seminar baru
	seminar := entity.Seminar{
		ID_Tiket:      idTiket,
		IDUser:        userID,
		PaymentStatus: "success",
		IsDeleted:     false,
	}

	if err := s.DB.Create(&seminar).Error; err != nil {
		return nil, err
	}

	return &dto.AdminAddParticipantResponse{
		Message:   "Berhasil menambahkan participant ke seminar",
		Status:    "success",
		IDTiket:   seminar.ID_Tiket,
		SeminarID: seminar.ID_Seminar,
		User: dto.UserInfo{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Jenjang:   user.Jenjang,
			Institusi: user.Institusi,
		},
	}, nil
}
