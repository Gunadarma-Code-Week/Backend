package service

import (
	"encoding/json"
	"errors"
	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
)

type TimelineService interface {
	CreateTimeline(dto dto.TimelineCreateDTO) (dto.TimelineResponseDTO, error)
	UpdateTimeline(id uint, dto dto.TimelineUpdateDTO) (dto.TimelineResponseDTO, error)
	DeleteTimeline(id uint) error
	GetTimelinesByCategory(category string) ([]dto.TimelineResponseDTO, error)
	GetTimelineByID(id uint) (dto.TimelineResponseDTO, error)
}

type timelineService struct {
	timelineRepository repository.TimelineRepository
}

func NewTimelineService(timelineRepo repository.TimelineRepository) TimelineService {
	return &timelineService{
		timelineRepository: timelineRepo,
	}
}

func (s *timelineService) CreateTimeline(d dto.TimelineCreateDTO) (dto.TimelineResponseDTO, error) {
	eventsJSON, _ := json.Marshal(d.Events)
	if d.Events == nil {
		eventsJSON = []byte("[]")
	}

	timeline := entity.Timeline{
		Category:    d.Category,
		OrderIndex:  d.OrderIndex,
		Date:        d.Date,
		Title:       d.Title,
		Description: d.Description,
		Events:      string(eventsJSON),
	}

	res, err := s.timelineRepository.InsertTimeline(timeline)
	if err != nil {
		return dto.TimelineResponseDTO{}, err
	}

	return s.toDTO(res), nil
}

func (s *timelineService) UpdateTimeline(id uint, d dto.TimelineUpdateDTO) (dto.TimelineResponseDTO, error) {
	timeline, err := s.timelineRepository.FindTimelineByID(id)
	if err != nil {
		return dto.TimelineResponseDTO{}, errors.New("timeline not found")
	}

	if d.Category != "" {
		timeline.Category = d.Category
	}
	if d.Date != "" {
		timeline.Date = d.Date
	}
	if d.Title != "" {
		timeline.Title = d.Title
	}
	timeline.OrderIndex = d.OrderIndex
	timeline.Description = d.Description

	if d.Events != nil {
		eventsJSON, _ := json.Marshal(d.Events)
		timeline.Events = string(eventsJSON)
	}

	res, err := s.timelineRepository.UpdateTimeline(timeline)
	if err != nil {
		return dto.TimelineResponseDTO{}, err
	}

	return s.toDTO(res), nil
}

func (s *timelineService) DeleteTimeline(id uint) error {
	timeline, err := s.timelineRepository.FindTimelineByID(id)
	if err != nil {
		return errors.New("timeline not found")
	}
	return s.timelineRepository.DeleteTimeline(timeline)
}

func (s *timelineService) GetTimelinesByCategory(category string) ([]dto.TimelineResponseDTO, error) {
	timelines, err := s.timelineRepository.AllTimelineByCategory(category)
	if err != nil {
		return nil, err
	}

	var res []dto.TimelineResponseDTO
	for _, t := range timelines {
		res = append(res, s.toDTO(t))
	}
	return res, nil
}

func (s *timelineService) GetTimelineByID(id uint) (dto.TimelineResponseDTO, error) {
	timeline, err := s.timelineRepository.FindTimelineByID(id)
	if err != nil {
		return dto.TimelineResponseDTO{}, errors.New("timeline not found")
	}
	return s.toDTO(timeline), nil
}

func (s *timelineService) toDTO(t entity.Timeline) dto.TimelineResponseDTO {
	var events []string
	if t.Events != "" && t.Events != "null" {
		json.Unmarshal([]byte(t.Events), &events)
	}
	if events == nil {
		events = []string{}
	}

	return dto.TimelineResponseDTO{
		ID:          t.ID,
		Category:    t.Category,
		OrderIndex:  t.OrderIndex,
		Date:        t.Date,
		Title:       t.Title,
		Description: t.Description,
		Events:      events,
	}
}
