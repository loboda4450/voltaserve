// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package repo

import (
	"encoding/json"
	"errors"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/log"
	"voltaserve/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SnapshotRepo interface {
	Find(id string) (model.Snapshot, error)
	FindByVersion(version int64) (model.Snapshot, error)
	FindAllForFile(fileID string) ([]model.Snapshot, error)
	FindAllDangling() ([]model.Snapshot, error)
	FindAllPrevious(fileID string, version int64) ([]model.Snapshot, error)
	GetIDsForFile(fileID string) ([]string, error)
	Insert(snapshot model.Snapshot) error
	Save(snapshot model.Snapshot) error
	Delete(id string) error
	Update(id string, opts SnapshotUpdateOptions) error
	MapWithFile(id string, fileID string) error
	DeleteMappingsForFile(fileID string) error
	DeleteAllDangling() error
	GetLatestVersionForFile(fileID string) (int64, error)
	CountAssociations(id string) (int, error)
	Attach(sourceFileID string, targetFileID string) error
	Detach(id string, fileID string) error
}

func NewSnapshotRepo() SnapshotRepo {
	return newSnapshotRepo()
}

func NewSnapshot() model.Snapshot {
	return &snapshotEntity{}
}

type snapshotEntity struct {
	ID         string         `json:"id" gorm:"column:id;size:36"`
	Version    int64          `json:"version" gorm:"column:version"`
	Original   datatypes.JSON `json:"original,omitempty" gorm:"column:original"`
	Preview    datatypes.JSON `json:"preview,omitempty" gorm:"column:preview"`
	Text       datatypes.JSON `json:"text,omitempty" gorm:"column:text"`
	OCR        datatypes.JSON `json:"ocr,omitempty" gorm:"column:ocr"`
	Entities   datatypes.JSON `json:"entities,omitempty" gorm:"column:entities"`
	Mosaic     datatypes.JSON `json:"mosaic,omitempty" gorm:"column:mosaic"`
	Watermark  datatypes.JSON `json:"watermark,omitempty" gorm:"column:watermark"`
	Thumbnail  datatypes.JSON `json:"thumbnail,omitempty" gorm:"column:thumbnail"`
	Status     string         `json:"status,omitempty" gorm:"column,status"`
	Language   *string        `json:"language,omitempty" gorm:"column:language"`
	TaskID     *string        `json:"taskID,omitempty" gorm:"column:task_id"`
	CreateTime string         `json:"createTime" gorm:"column:create_time"`
	UpdateTime *string        `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (*snapshotEntity) TableName() string {
	return "snapshot"
}

func (s *snapshotEntity) BeforeCreate(*gorm.DB) (err error) {
	s.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (s *snapshotEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	s.UpdateTime = &timeNow
	return nil
}

func (s *snapshotEntity) GetID() string {
	return s.ID
}

func (s *snapshotEntity) GetVersion() int64 {
	return s.Version
}

func (s *snapshotEntity) GetOriginal() *model.S3Object {
	if s.Original.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Original.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetPreview() *model.S3Object {
	if s.Preview.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Preview.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetText() *model.S3Object {
	if s.Text.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Text.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetOCR() *model.S3Object {
	if s.OCR.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.OCR.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetEntities() *model.S3Object {
	if s.Entities.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Entities.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetMosaic() *model.S3Object {
	if s.Mosaic.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Mosaic.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetWatermark() *model.S3Object {
	if s.Watermark.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Watermark.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetThumbnail() *model.S3Object {
	if s.Thumbnail.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Thumbnail.String()), &res); err != nil {
		log.GetLogger().Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) GetStatus() string {
	return s.Status
}

func (s *snapshotEntity) GetLanguage() *string {
	return s.Language
}

func (s *snapshotEntity) GetTaskID() *string {
	return s.TaskID
}

func (s *snapshotEntity) SetID(id string) {
	s.ID = id
}

func (s *snapshotEntity) SetVersion(version int64) {
	s.Version = version
}

func (s *snapshotEntity) SetOriginal(m *model.S3Object) {
	if m == nil {
		s.Original = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Original.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetPreview(m *model.S3Object) {
	if m == nil {
		s.Preview = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Preview.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetText(m *model.S3Object) {
	if m == nil {
		s.Text = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Text.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetOCR(m *model.S3Object) {
	if m == nil {
		s.OCR = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.OCR.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetEntities(m *model.S3Object) {
	if m == nil {
		s.Entities = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Entities.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetMosaic(m *model.S3Object) {
	if m == nil {
		s.Mosaic = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Mosaic.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetWatermark(m *model.S3Object) {
	if m == nil {
		s.Watermark = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Watermark.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetThumbnail(m *model.S3Object) {
	if m == nil {
		s.Thumbnail = nil
	} else {
		b, err := json.Marshal(m)
		if err != nil {
			log.GetLogger().Fatal(err)
			return
		}
		if err := s.Thumbnail.UnmarshalJSON(b); err != nil {
			log.GetLogger().Fatal(err)
		}
	}
}

func (s *snapshotEntity) SetStatus(status string) {
	s.Status = status
}

func (s *snapshotEntity) SetLanguage(language string) {
	s.Language = &language
}

func (s *snapshotEntity) SetTaskID(taskID *string) {
	s.TaskID = taskID
}

func (s *snapshotEntity) HasOriginal() bool {
	return s.Original != nil
}

func (s *snapshotEntity) HasPreview() bool {
	return s.Preview != nil
}

func (s *snapshotEntity) HasText() bool {
	return s.Text != nil
}

func (s *snapshotEntity) HasOCR() bool {
	return s.OCR != nil
}

func (s *snapshotEntity) HasEntities() bool {
	return s.Entities != nil
}

func (s *snapshotEntity) HasMosaic() bool {
	return s.Mosaic != nil
}

func (s *snapshotEntity) HasWatermark() bool {
	return s.Watermark != nil
}

func (s *snapshotEntity) HasThumbnail() bool {
	return s.Thumbnail != nil
}

func (s *snapshotEntity) GetCreateTime() string {
	return s.CreateTime
}

func (s *snapshotEntity) GetUpdateTime() *string {
	return s.UpdateTime
}

type snapshotRepo struct {
	db *gorm.DB
}

func newSnapshotRepo() *snapshotRepo {
	return &snapshotRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *snapshotRepo) find(id string) (*snapshotEntity, error) {
	var res snapshotEntity
	if db := repo.db.Where("id = ?", id).First(&res); db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewSnapshotNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *snapshotRepo) Find(id string) (model.Snapshot, error) {
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *snapshotRepo) FindByVersion(version int64) (model.Snapshot, error) {
	var res = snapshotEntity{}
	db := repo.db.Where("version = ?", version).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewSnapshotNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *snapshotRepo) Insert(snapshot model.Snapshot) error {
	if db := repo.db.Create(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) Save(snapshot model.Snapshot) error {
	if db := repo.db.Save(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) Delete(id string) error {
	snapshot, err := repo.find(id)
	if err != nil {
		return err
	}
	if db := repo.db.Delete(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

type SnapshotUpdateOptions struct {
	Fields    []string `json:"fields"`
	Original  *model.S3Object
	Preview   *model.S3Object
	Text      *model.S3Object
	OCR       *model.S3Object
	Entities  *model.S3Object
	Mosaic    *model.S3Object
	Watermark *model.S3Object
	Thumbnail *model.S3Object
	Status    *string
	Language  *string
	TaskID    *string
}

const (
	SnapshotFieldOriginal  = "original"
	SnapshotFieldPreview   = "preview"
	SnapshotFieldText      = "text"
	SnapshotFieldOCR       = "ocr"
	SnapshotFieldEntities  = "entities"
	SnapshotFieldMosaic    = "mosaic"
	SnapshotFieldWatermark = "watermark"
	SnapshotFieldThumbnail = "thumbnail"
	SnapshotFieldStatus    = "status"
	SnapshotFieldLanguage  = "language"
	SnapshotFieldTaskID    = "taskID"
)

func (repo *snapshotRepo) Update(id string, opts SnapshotUpdateOptions) error {
	snapshot, err := repo.find(id)
	if err != nil {
		return err
	}
	if helper.Includes(opts.Fields, SnapshotFieldOriginal) {
		snapshot.SetOriginal(opts.Original)
	}
	if helper.Includes(opts.Fields, SnapshotFieldPreview) {
		snapshot.SetPreview(opts.Preview)
	}
	if helper.Includes(opts.Fields, SnapshotFieldText) {
		snapshot.SetText(opts.Text)
	}
	if helper.Includes(opts.Fields, SnapshotFieldOCR) {
		snapshot.SetOCR(opts.OCR)
	}
	if helper.Includes(opts.Fields, SnapshotFieldEntities) {
		snapshot.SetEntities(opts.Entities)
	}
	if helper.Includes(opts.Fields, SnapshotFieldMosaic) {
		snapshot.SetMosaic(opts.Mosaic)
	}
	if helper.Includes(opts.Fields, SnapshotFieldWatermark) {
		snapshot.SetWatermark(opts.Watermark)
	}
	if helper.Includes(opts.Fields, SnapshotFieldThumbnail) {
		snapshot.SetThumbnail(opts.Thumbnail)
	}
	if helper.Includes(opts.Fields, SnapshotFieldStatus) {
		snapshot.SetStatus(*opts.Status)
	}
	if helper.Includes(opts.Fields, SnapshotFieldLanguage) {
		snapshot.SetLanguage(*opts.Language)
	}
	if helper.Includes(opts.Fields, SnapshotFieldTaskID) {
		snapshot.SetTaskID(opts.TaskID)
	}
	if db := repo.db.Save(&snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) MapWithFile(id string, fileID string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id) VALUES (?, ?)", id, fileID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) DeleteMappingsForFile(fileID string) error {
	if db := repo.db.Exec("DELETE FROM snapshot_file WHERE file_id = ?", fileID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) findAllForFile(fileID string) ([]*snapshotEntity, error) {
	var res []*snapshotEntity
	db := repo.db.
		Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.file_id = ? ORDER BY s.version", fileID).
		Scan(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return res, nil
}

func (repo *snapshotRepo) FindAllForFile(fileID string) ([]model.Snapshot, error) {
	snapshots, err := repo.findAllForFile(fileID)
	if err != nil {
		return nil, err
	}
	var res []model.Snapshot
	for _, s := range snapshots {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) FindAllDangling() ([]model.Snapshot, error) {
	var entities []*snapshotEntity
	db := repo.db.Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL").Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Snapshot
	for _, s := range entities {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) FindAllPrevious(fileID string, version int64) ([]model.Snapshot, error) {
	var entities []*snapshotEntity
	db := repo.db.Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.file_id = ? AND s.version < ? ORDER BY s.version DESC", fileID, version).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Snapshot
	for _, s := range entities {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) GetIDsForFile(fileID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw("SELECT snapshot_id result FROM snapshot_file WHERE file_id = ?", fileID).
		Scan(&values)
	if db.Error != nil {
		return nil, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *snapshotRepo) DeleteAllDangling() error {
	if db := repo.db.Exec("DELETE FROM snapshot WHERE id IN (SELECT s.id FROM (SELECT * FROM snapshot) s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL)"); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) GetLatestVersionForFile(fileID string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	if db := repo.db.
		Raw("SELECT coalesce(max(s.version), 0) result FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id WHERE map.file_id = ?", fileID).
		Scan(&res); db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}

func (repo *snapshotRepo) CountAssociations(id string) (int, error) {
	type Result struct {
		Count int
	}
	var res Result
	if db := repo.db.Raw("SELECT COUNT(*) count FROM snapshot_file WHERE snapshot_id = ?", id).Scan(&res); db.Error != nil {
		return 0, db.Error
	}
	return res.Count, nil
}

func (repo *snapshotRepo) Attach(sourceFileID string, targetFileID string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id) SELECT s.id, ? "+
		"FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id "+
		"WHERE map.file_id = ? ORDER BY s.version DESC LIMIT 1", targetFileID, sourceFileID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) Detach(id string, fileID string) error {
	if db := repo.db.Exec("DELETE FROM snapshot_file WHERE snapshot_id = ? AND file_id = ?", id, fileID); db.Error != nil {
		return db.Error
	}
	return nil
}
