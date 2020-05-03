package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// NotificationsRepo is used to contain the db driver.
type NotificationsRepo struct {
    DB *pg.DB
}

func (m *NotificationsRepo) CreateMany(notifications []models.Notification) ([]models.Notification, error) {
    _, err := m.DB.Model(&notifications).Insert()
    return notifications, err
}

// Create is used to create notification for the database.
func (m *NotificationsRepo) Create(notification *models.Notification) (*models.Notification, error) {
    _, err := m.DB.Model(notification).Returning("*").Insert()
    return notification, err
}

// Update is used to update the passed notification by id.
func (m *NotificationsRepo) Update(notification *models.Notification) (*models.Notification, error) {
    _, err := m.DB.Model(notification).Where("id = ?", notification.ID).Update()
    return notification, err
}

// GetByID is used to fetch notification by id.
func (m *NotificationsRepo) GetByID(id string) (*models.Notification, error) {
    notification := models.Notification{}
    err := m.DB.Model(&notification).Where("id = ?", id).Select()
    if err != nil {
        return nil, err
    }
    return &notification, nil
}

// Delete is used to delete notification by its id.
func (m *NotificationsRepo) Delete(notification *models.Notification) error {
    _, err := m.DB.Model(notification).Where("id = ?", notification.ID).Delete()
    return err
}
