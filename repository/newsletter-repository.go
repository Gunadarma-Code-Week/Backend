package repository

import (
	"gcw/entity"
	"gorm.io/gorm"
)

type NewsletterRepository interface {
	FindById(uint64) (*entity.NewsLetter, error)
	Create(*entity.NewsLetter) error
	Update(*entity.NewsLetter) error
	Delete(*entity.NewsLetter) error
}

type newsletterRepository struct {
	DB *gorm.DB
}

func NewNewsletterRepository(db *gorm.DB) NewsletterRepository {
	return &newsletterRepository{DB: db}
}

func (r *newsletterRepository) FindById(id uint64) (*entity.NewsLetter, error) {
	var n entity.NewsLetter
	if err := r.DB.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *newsletterRepository) Create(n *entity.NewsLetter) error {
	return r.DB.Create(n).Error
}

func (r *newsletterRepository) Update(n *entity.NewsLetter) error {
	return r.DB.Save(n).Error
}

func (r *newsletterRepository) Delete(n *entity.NewsLetter) error {
	return r.DB.Delete(n).Error
}
