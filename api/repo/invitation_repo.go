package repo

import (
	"errors"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type InvitationInsertOptions struct {
	UserId         string
	OrganizationId string
	Emails         []string
}

type InvitationRepo interface {
	Insert(opts InvitationInsertOptions) ([]model.Invitation, error)
	Find(id string) (model.Invitation, error)
	GetIncoming(email string) ([]model.Invitation, error)
	GetOutgoing(organizationId string, userId string) ([]model.Invitation, error)
	Save(org model.Invitation) error
	Delete(id string) error
}

func NewInvitationRepo() InvitationRepo {
	return newInvitationRepo()
}

func NewInvitation() model.Invitation {
	return &invitationEntity{}
}

type invitationEntity struct {
	ID             string  `json:"id"`
	OrganizationId string  `json:"organizationId"`
	OwnerId        string  `json:"ownerId"`
	Email          string  `json:"email"`
	Status         string  `json:"status"`
	CreateTime     string  `json:"createTime"`
	UpdateTime     *string `json:"updateTime"`
}

func (invitationEntity) TableName() string {
	return "invitation"
}

func (o *invitationEntity) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (o *invitationEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	o.UpdateTime = &timeNow
	return nil
}

func (i invitationEntity) GetID() string {
	return i.ID
}

func (i invitationEntity) GetOrganizationID() string {
	return i.OrganizationId
}

func (i invitationEntity) GetOwnerID() string {
	return i.OwnerId
}

func (i invitationEntity) GetEmail() string {
	return i.Email
}

func (i invitationEntity) GetStatus() string {
	return i.Status
}

func (i invitationEntity) GetCreateTime() string {
	return i.CreateTime
}

func (i invitationEntity) GetUpdateTime() *string {
	return i.UpdateTime
}

func (w *invitationEntity) SetStatus(status string) {
	w.Status = status
}

func (w *invitationEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type invitationRepo struct {
	db       *gorm.DB
	userRepo *userRepo
}

func newInvitationRepo() *invitationRepo {
	return &invitationRepo{
		db:       infra.GetDb(),
		userRepo: newUserRepo(),
	}
}

func (repo *invitationRepo) Insert(opts InvitationInsertOptions) ([]model.Invitation, error) {
	var res []model.Invitation
	for _, e := range opts.Emails {
		invitation := invitationEntity{
			ID:             helpers.NewId(),
			OrganizationId: opts.OrganizationId,
			OwnerId:        opts.UserId,
			Email:          e,
			Status:         model.InvitationStatusPending,
		}
		if db := repo.db.Save(&invitation); db.Error != nil {
			return nil, db.Error
		}
		i, err := repo.Find(invitation.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func (repo *invitationRepo) Find(id string) (model.Invitation, error) {
	var invitation = invitationEntity{}
	db := repo.db.Where("id = ?", id).First(&invitation)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewInvitationNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &invitation, nil
}

func (repo *invitationRepo) GetIncoming(email string) ([]model.Invitation, error) {
	var invitations []*invitationEntity
	db := repo.db.
		Raw("SELECT * FROM invitation WHERE email = ? and status = 'pending' ORDER BY create_time DESC", email).
		Scan(&invitations)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Invitation
	for _, inv := range invitations {
		res = append(res, inv)
	}
	return res, nil
}

func (repo *invitationRepo) GetOutgoing(organizationId string, userId string) ([]model.Invitation, error) {
	var invitations []*invitationEntity
	db := repo.db.
		Raw("SELECT * FROM invitation WHERE organization_id = ? and owner_id = ? ORDER BY create_time DESC", organizationId, userId).
		Scan(&invitations)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Invitation
	for _, inv := range invitations {
		res = append(res, inv)
	}
	return res, nil
}

func (repo *invitationRepo) Save(org model.Invitation) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *invitationRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM invitation WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
