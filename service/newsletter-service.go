package service

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
	"github.com/mashingan/smapping"
)

type NewsletterService interface {
	Create(dto.CreateNewsLetterDTO) (*entity.NewsLetter, error)
	FindByID(uint64) (*entity.NewsLetter, error)
	Update(uint64, dto.UpdateNewsLetterDTO) (*entity.NewsLetter, error)
	Delete(uint64) error
}

type newsletterService struct {
	repo repository.NewsletterRepository
}

func NewNewsletterService(r repository.NewsletterRepository) NewsletterService {
	return &newsletterService{repo: r}
}

func (s *newsletterService) Create(input dto.CreateNewsLetterDTO) (*entity.NewsLetter, error) {
	newsletter := entity.NewsLetter{}
	smapping.FillStruct(&newsletter, smapping.MapFields(input))
	err := s.repo.Create(&newsletter)
	return &newsletter, err
}

func (s *newsletterService) FindByID(id uint64) (*entity.NewsLetter, error) {
	return s.repo.FindById(id)
}

func (s *newsletterService) Update(id uint64, input dto.UpdateNewsLetterDTO) (*entity.NewsLetter, error) {
	newsletter, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	smapping.FillStruct(newsletter, smapping.MapFields(input))
	err = s.repo.Update(newsletter)
	return newsletter, err
}

func (s *newsletterService) Delete(id uint64) error {
	newsletter, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(newsletter)
}
