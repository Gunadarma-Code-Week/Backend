package repository

import (
	"gcw/entity"

	"gorm.io/gorm"
)

type TimelineRepository interface {
	InsertTimeline(timeline entity.Timeline) (entity.Timeline, error)
	UpdateTimeline(timeline entity.Timeline) (entity.Timeline, error)
	DeleteTimeline(timeline entity.Timeline) error
	AllTimelineByCategory(category string) ([]entity.Timeline, error)
	FindTimelineByID(timelineID uint) (entity.Timeline, error)
}

type timelineConnection struct {
	connection *gorm.DB
}

func NewTimelineRepository(dbConn *gorm.DB) TimelineRepository {
	return &timelineConnection{
		connection: dbConn,
	}
}

func (db *timelineConnection) InsertTimeline(timeline entity.Timeline) (entity.Timeline, error) {
	err := db.connection.Save(&timeline).Error
	return timeline, err
}

func (db *timelineConnection) UpdateTimeline(timeline entity.Timeline) (entity.Timeline, error) {
	err := db.connection.Save(&timeline).Error
	return timeline, err
}

func (db *timelineConnection) DeleteTimeline(timeline entity.Timeline) error {
	return db.connection.Delete(&timeline).Error
}

func (db *timelineConnection) AllTimelineByCategory(category string) ([]entity.Timeline, error) {
	var timelines []entity.Timeline
	err := db.connection.Where("category = ?", category).Order("order_index ASC").Find(&timelines).Error
	return timelines, err
}

func (db *timelineConnection) FindTimelineByID(timelineID uint) (entity.Timeline, error) {
	var timeline entity.Timeline
	err := db.connection.Where("id = ?", timelineID).Take(&timeline).Error
	return timeline, err
}
