package resolvers

import (
    "context"
    "fmt"
    "log"

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
    // TODO: Delete members of group, on deletion of group.
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

func (m *mutationResolver) notifyMembersOfGroupForMeetupCreation(group *models.Group, meetup *models.Meetup) {
    notificationsLength := len(group.Members)
    if group.UserID != meetup.UserID {
        // prepare the length of the notificaion to be capable of contain the extra index for the gorup owner
        // which doesn't exist at the group.Members
        notificationsLength++
    }
    notifications := make([]models.Notification, notificationsLength)
    for i := 0; i < len(group.Members); i++ {
        notification := models.Notification{
            UserID:         group.Members[i].ID, // the one we wish to notify which is the member of group.
            NotifiableType: "meetup_created",
            NotifiableID:   meetup.ID, // the meetup that caused the notification.
        }
        notifications[i] = notification
    }
    if group.UserID != meetup.UserID {
        // notify the owner of the group as well.
        notifications[notificationsLength-1] = models.Notification{
            UserID:         group.UserID,
            NotifiableType: "meetup_created",
            NotifiableID:   meetup.ID,
        }
    }
    // we need to notify
    notifications, err := m.NotificationsRepo.CreateMany(notifications)
    if err != nil {
        log.Fatal(err)
    }
    for _, notification := range notifications {
        m.nClient.Publish("notification.user_"+notification.UserID, &notification)
    }

}
func (m *mutationResolver) CreateGroupMeetup(ctx context.Context, input models.CreateMeetupInput, groupID string) (*models.Meetup, error) {

    meetup, err := m.CreateMeetup(ctx, input)
    if err != nil {
        return nil, err
    }
    status, err := m.ShareMeetupToGroup(ctx, meetup.ID, groupID)
    if err != nil || !status {
        return nil, errors.ErrInternalError
    }
    group, err := m.GroupsRepo.GetGroupMembersExceptFor(groupID, meetup.UserID)
    if err != nil {
        return nil, err
    }
    go m.notifyMembersOfGroupForMeetupCreation(group, meetup)

    return meetup, nil
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

func (m *mutationResolver) LeaveGroup(ctx context.Context, id string) (*models.Group, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    group, err := m.GroupsRepo.GetByID(id)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }

    exists, err := m.GroupsRepo.IsUserMemberOfGroup(id, currentUser.ID)
    if err != nil || !exists {
        return nil, errors.ErrRecordNotFound
    }
    return m.GroupsRepo.DischargeMemberFromGroup(group, currentUser.ID)

}
func (m *mutationResolver) DischargeMemberFromGroup(ctx context.Context, id string, userID string) (*models.Group, error) {
    currentUser, _ := middlewares.GetCurrentUserFromContext(ctx)
    group, err := m.GroupsRepo.GetByID(id)
    if err != nil {
        return nil, errors.ErrRecordNotFound
    }
    // if he isn't an creator of the group
    // or the currently authenticated user is not a secondary admin nor a moderator.
    exists, err := m.GroupsRepo.IsUserSecondaryAdminOfGroup(id, currentUser.ID)
    if err != nil || (!exists && group.UserID != currentUser.ID) {
        return nil, errors.ErrUnauthenticated
    }
    return m.GroupsRepo.DischargeMemberFromGroup(group, userID)
}
