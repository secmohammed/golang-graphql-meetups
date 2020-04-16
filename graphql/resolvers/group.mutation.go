package resolvers

import (
    "context"
    "fmt"

    "github.com/secmohammed/meetups/middlewares"
    "github.com/secmohammed/meetups/models"
    "github.com/secmohammed/meetups/utils/errors"
)

func (m *mutationResolver) CreateGroup(ctx context.Context, input models.CreateGroupInput) (*models.Group, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    if err := input.Validate(); err != nil {
        return nil, err
    }
    group := &models.Group{
        Name:        input.Name,
        Description: input.Description,
        UserID:      currentUser.ID,
    }
    // TODO: attach user_id to group_id, attach categories_ids to group_id
    return m.GroupsRepo.Create(group)
}
func (m *mutationResolver) DeleteGroup(ctx context.Context, id string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)

    group, err := m.GroupsRepo.GetByID(id)
    if err != nil || group == nil {
        return false, errors.ErrRecordNotFound
    }
    if group.UserID != currentUser.ID {
        return false, errors.ErrUnauthenticated
    }

    err = m.GroupsRepo.Delete(group)
    if err != nil {
        return false, fmt.Errorf("error while deleting group: %v", err)
    }
    return true, nil

}
func (m *mutationResolver) UpdateGroup(ctx context.Context, id string, input models.UpdateGroupInput) (*models.Group, error) {
    return nil, nil
}
func (m *mutationResolver) AssignMemberToGroup(ctx context.Context, id string, userID string) (*models.Group, error) {
    return nil, nil
}
