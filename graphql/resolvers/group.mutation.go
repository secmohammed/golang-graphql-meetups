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

    group, err := m.GroupsRepo.Create(group)
    if err != nil {
        return nil, err
    }
    m.GroupsRepo.AttachCategoriesToGroup(input.CategoryIds, group)
    return group, nil
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
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    if err := input.Validate(); err != nil {
        return nil, err
    }
    group, err := m.GroupsRepo.GetByID(id)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    // if he is an creator of the group
    // or currently authenticated user is a secondary admin or a moderator.
    exists, err := m.GroupsRepo.IsUserSecondaryAdminOfGroup(id, currentUser.ID)
    if err != nil || (!exists && group.UserID != currentUser.ID) {
        return nil, errors.ErrCouldntAssignMemberToGroup
    }
    if input.Name != "" {
        group.Name = input.Name
    }
    if input.Description != "" {
        group.Description = input.Description
    }
    if len(input.CategoryIds) != 0 {
        err := m.GroupsRepo.SyncCategoriesWithGroup(input.CategoryIds, group)
        if err != nil {
            return nil, err
        }
    }
    return m.GroupsRepo.Update(group)
}
func (m *mutationResolver) DeleteMeetupFromGroup(ctx context.Context, meetupID string, groupID string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    group, err := m.GroupsRepo.GetByID(groupID)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    // if he is an creator of the group
    // or currently authenticated user is a secondary admin or a moderator.
    exists, err := m.GroupsRepo.IsUserSecondaryAdminOfGroup(groupID, currentUser.ID)
    if err != nil || (!exists && group.UserID != currentUser.ID) {
        return false, errors.ErrCouldntAssignMemberToGroup
    }

    return true, m.GroupsRepo.DetachMeetupFromGroup(groupID, meetupID)

}
func (m *mutationResolver) ShareMeetupToGroup(ctx context.Context, meetupID string, groupID string) (bool, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    group, err := m.GroupsRepo.GetByID(groupID)
    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    // make sure that the current authenticated user is member/owner of the group.
    exists, err := m.GroupsRepo.IsUserMemberOfGroup(groupID, currentUser.ID)
    if err != nil || (!exists && group.UserID != currentUser.ID) {
        return false, errors.ErrCouldntAssignMemberToGroup
    }

    meetup, err := m.MeetupsRepo.GetByID(meetupID)

    if err != nil {
        return false, errors.ErrRecordNotFound
    }
    return true, m.GroupsRepo.AttachMeetupToGroup(group, meetup)
}
func (m *mutationResolver) AssignMemberToGroup(ctx context.Context, id string, userID string, role *string) (*models.Group, error) {
    // given we have the group id and the user id we want to assign to this group
    // when the authenticated user has the ability to add others based on
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    group, err := m.GroupsRepo.GetByID(id)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    // if he is an creator of the group
    // or currently authenticated user is a secondary admin or a moderator.
    exists, err := m.GroupsRepo.IsUserSecondaryAdminOfGroup(id, currentUser.ID)
    if err != nil || (!exists && group.UserID != currentUser.ID) {
        return nil, errors.ErrCouldntAssignMemberToGroup
    }
    // then we can add user to group.

    return m.GroupsRepo.AssignMemberToGroup(group, userID, *role)
}
