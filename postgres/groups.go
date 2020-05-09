package postgres

import (
    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
    "github.com/secmohammed/meetups/models"
)

// GroupsRepo is used to contain the db driver.
type GroupsRepo struct {
    DB *pg.DB
}

func (g *GroupsRepo) DetachMeetupFromGroup(groupID, meetupID string) error {
    meetupGroup := models.MeetupGroup{
        MeetupID: meetupID,
        GroupID:  groupID,
    }
    _, err := g.DB.Model(&meetupGroup).Where("meetup_id = ?", meetupID).Where("group_id = ?", groupID).Delete()
    return err

}
func (g *GroupsRepo) AttachMeetupToGroup(group *models.Group, meetup *models.Meetup) error {
    meetupGroup := models.MeetupGroup{
        MeetupID: meetup.ID,
        GroupID:  group.ID,
    }
    _, err := g.DB.Model(&meetupGroup).Insert()
    return err
}
func (g *GroupsRepo) DischargeMemberFromGroup(group *models.Group, userID string) (*models.Group, error) {
    groupUser := models.GroupUser{
        GroupID: group.ID,
        UserID:  userID,
    }
    _, err := g.DB.Model(&groupUser).Where("group_id = ?", group.ID).Where("user_id = ? ", userID).Delete()
    if err != nil {
        return nil, err
    }
    return group, nil

}
func (g *GroupsRepo) AssignMemberToGroup(group *models.Group, userID, role string) (*models.Group, error) {
    //TODO: FirstOrCreate.
    groupUser := models.GroupUser{
        GroupID: group.ID,
        UserID:  userID,
        Type:    role,
    }
    _, err := g.DB.Model(&groupUser).Insert()
    if err != nil {
        return nil, err
    }
    return group, nil
}

func (g *GroupsRepo) SyncCategoriesWithGroup(categoryIds []string, group *models.Group) error {
    categoryGroup := models.CategoryGroup{
        GroupID: group.ID,
    }
    _, err := g.DB.Model(&categoryGroup).Where("group_id = ?", group.ID).Delete()
    if err != nil {
        return err
    }
    return g.AttachCategoriesToGroup(categoryIds, group)

}
func (g *GroupsRepo) AttachCategoriesToGroup(categoryIds []string, group *models.Group) error {
    categories := make([]models.CategoryGroup, len(categoryIds))
    for i := 0; i < len(categoryIds); i++ {
        categories = append(categories, models.CategoryGroup{
            GroupID:    group.ID,
            CategoryID: categoryIds[i],
        })

    }

    _, err := g.DB.Model(&categories).Insert()
    return err
}

//Create is used to create a group using the passed struct.
func (g *GroupsRepo) Create(group *models.Group) (*models.Group, error) {
    _, err := g.DB.Model(group).Returning("*").Insert()
    return group, err
}
func (g *GroupsRepo) GetGroupMembersExceptFor(groupID, userID string) (*models.Group, error) {
    group := models.Group{}
    err := g.DB.Model(&group).Where("\"group\".\"id\" = ?", groupID).Relation("Members", func(q *orm.Query) (*orm.Query, error) {
        return q.Where("group_user.user_id != ?", userID), nil
    }).First()
    if err != nil {
        return nil, err
    }
    return &group, nil

}
func (g *GroupsRepo) GetGroupMembers(groupID string) (*models.Group, error) {
    group := models.Group{}
    err := g.DB.Model(&group).Where("\"group\".\"id\" = ?", groupID).Relation("Members").First()
    if err != nil {
        return nil, err
    }
    return &group, nil

}
func (g *GroupsRepo) IsUserMemberOfGroup(id, authenticated_user_id string) (bool, error) {
    group := models.Group{}
    err := g.DB.Model(&group).Where("\"group\".\"id\" = ?", id).Relation("Members", func(q *orm.Query) (*orm.Query, error) {
        return q.Where("group_user.user_id = ?", authenticated_user_id), nil
    }).First()
    boolValue := !(len(group.Members) == 0)
    if err != nil {
        return boolValue, err
    }
    return boolValue, nil

}
func (g *GroupsRepo) IsUserSecondaryAdminOfGroup(id, authenticated_user_id string) (bool, error) {
    group := models.Group{}
    err := g.DB.Model(&group).Where("\"group\".\"id\" = ?", id).Relation("Members", func(q *orm.Query) (*orm.Query, error) {
        return q.Where("group_user.user_id = ?", authenticated_user_id).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
            return q.WhereOr("group_user.type = ?", "moderator").WhereOr("group_user.type = ?", "admin"), nil
        }), nil
    }).First()
    boolValue := !(len(group.Members) == 0)
    if err != nil {
        return boolValue, err
    }
    return boolValue, nil
}

// GetByID is used to fetch meetup by id.
func (g *GroupsRepo) GetByID(id string) (*models.Group, error) {
    group := models.Group{}
    // err := g.DB.Model(&group).Where("\"group\".\"id\" = ?", id).Relation("Members").First()
    err := g.DB.Model(&group).Where("id = ?", id).First()
    if err != nil {
        return nil, err
    }
    return &group, nil
}

// GetGroups is used to get meetups from database.
func (g *GroupsRepo) GetGroups() ([]*models.Group, error) {
    var groups []*models.Group
    err := g.DB.Model(&groups).Order("id").Select()
    if err != nil {
        return nil, err
    }
    return groups, nil
}

// Update is used to update the passed group by id.
func (g *GroupsRepo) Update(group *models.Group) (*models.Group, error) {
    _, err := g.DB.Model(group).Where("id = ?", group.ID).Update()
    return group, err
}

// Delete is used to delete group by its id.
func (g *GroupsRepo) Delete(group *models.Group) error {
    _, err := g.DB.Model(group).Where("id = ?", group.ID).Delete()
    return err
}
