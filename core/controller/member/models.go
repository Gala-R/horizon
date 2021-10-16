package member

import (
	"context"
	"errors"
	"time"

	"g.hz.netease.com/horizon/pkg/member/models"
	memberservice "g.hz.netease.com/horizon/pkg/member/service"
	userManager "g.hz.netease.com/horizon/pkg/user/manager"
	usermodels "g.hz.netease.com/horizon/pkg/user/models"
)

var (
	// Converter us tge global member converter
	Converter = New()
)

var (
	Owner string = "Owner"
)

type PostMember struct {
	// ResourceType group/application/applicationInstance
	ResourceType string

	// ResourceID group id;application id ...
	ResourceID uint

	// MemberType user or group
	MemberType models.MemberType

	// MemberNameID group id / userid
	MemberNameID uint

	// Role owner/maintainer/develop/...
	Role string
}

type Member struct {
	// ID the uniq id of the member entry
	ID uint

	// ResourceName   application/group
	ResourceType models.ResourceType
	ResourceID   uint

	// MemberType user or group
	MemberType models.MemberType

	// MemberName username or groupName
	MemberName string
	// MemberNameID userID or groupID
	MemberNameID uint

	// Role the role name that bind
	Role string
	// GrantBy user who grant the role
	GrantBy uint
	// GrantTime
	GrantTime time.Time
}

func CovertPostMember(member *PostMember) memberservice.PostMember {
	return memberservice.PostMember{
		ResourceType: member.ResourceType,
		ResourceID:   member.ResourceID,
		MemberInfo:   member.MemberNameID,
		MemberType:   member.MemberType,
		Role:         member.Role,
	}
}

type ConvertMemberHelp interface {
	ConvertMember(ctx context.Context, member *models.Member) (*Member, error)
	ConvertMembers(ctx context.Context, member []models.Member) ([]Member, error)
}

type converter struct {
	userManager userManager.Manager
}

func New() ConvertMemberHelp {
	return &converter{
		userManager: userManager.Mgr,
	}
}

func (c *converter) ConvertMember(ctx context.Context, member *models.Member) (_ *Member, err error) {
	// convert userID to userName
	var memberInfo string
	var user *usermodels.User

	if member.MemberType == models.MemberUser {
		user, err = c.userManager.GetUserByID(ctx, member.MemberNameID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}
		memberInfo = user.Name
	} else {
		// TODO(tom) covert groupID to GroupName
		return nil, errors.New("group member not support yet")
	}

	return &Member{
		ID:           member.ID,
		MemberType:   member.MemberType,
		MemberName:   memberInfo,
		MemberNameID: member.MemberNameID,
		ResourceType: member.ResourceType,
		ResourceID:   member.ResourceID,
		Role:         member.Role,
		GrantBy:      member.GrantBy,
		GrantTime:    member.UpdatedAt,
	}, nil
}
func (c *converter) ConvertMembers(ctx context.Context, members []models.Member) ([]Member, error) {
	var userIDs []uint
	for _, member := range members {
		if member.MemberType != models.MemberUser {
			return nil, errors.New("Only Support User MemberType yet")
		}
		userIDs = append(userIDs, member.MemberNameID)
	}
	users, err := c.userManager.GetUserByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if len(users) != len(userIDs) {
		return nil, errors.New("cannot find all the users")
	}
	var retMembers []Member
	for i, member := range members {
		retMembers = append(retMembers, Member{
			ID:           member.ID,
			MemberType:   member.MemberType,
			MemberName:   users[i].Name,
			MemberNameID: member.MemberNameID,
			ResourceType: member.ResourceType,
			ResourceID:   member.ResourceID,
			Role:         member.Role,
			GrantBy:      member.GrantBy,
			GrantTime:    member.UpdatedAt,
		})
	}
	return retMembers, nil
}
