// Copyright © 2023 Horizoncd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"context"

	"github.com/horizoncd/horizon/lib/q"
	amodels "github.com/horizoncd/horizon/pkg/application/models"
	cmodels "github.com/horizoncd/horizon/pkg/cluster/models"
	"github.com/horizoncd/horizon/pkg/template/dao"
	"github.com/horizoncd/horizon/pkg/template/models"
	"gorm.io/gorm"
)

// nolint
//
//go:generate mockgen -source=$GOFILE -destination=../../../mock/pkg/template/manager/manager_mock.go -package=mock_manager
type Manager interface {
	// Create template
	Create(ctx context.Context, template *models.Template) (*models.Template, error)
	ListV2(ctx context.Context, query *q.Query, groupIDs ...uint) ([]*models.Template, error)
	// ListTemplate returns all template
	ListTemplate(ctx context.Context) ([]*models.Template, error)
	// ListByGroupID lists all template by group ID
	ListByGroupID(ctx context.Context, groupID uint) ([]*models.Template, error)
	// DeleteByID deletes template by ID
	DeleteByID(ctx context.Context, id uint) error
	// GetByID gets a template by ID
	GetByID(ctx context.Context, id uint) (*models.Template, error)
	// GetByName gets a template by name
	GetByName(ctx context.Context, name string) (*models.Template, error)
	GetRefOfApplication(ctx context.Context, id uint) ([]*amodels.Application, uint, error)
	GetRefOfCluster(ctx context.Context, id uint) ([]*cmodels.Cluster, uint, error)
	UpdateByID(ctx context.Context, id uint, template *models.Template) error
	ListByGroupIDs(ctx context.Context, ids []uint) ([]*models.Template, error)
	ListByIDs(ctx context.Context, ids []uint) ([]*models.Template, error)
}

func New(db *gorm.DB) Manager {
	return &manager{dao: dao.NewDAO(db)}
}

type manager struct {
	dao dao.DAO
}

func (m *manager) Create(ctx context.Context, template *models.Template) (*models.Template, error) {
	return m.dao.Create(ctx, template)
}

func (m *manager) ListTemplate(ctx context.Context) ([]*models.Template, error) {
	return m.dao.ListTemplate(ctx)
}

func (m *manager) ListByGroupID(ctx context.Context, groupID uint) ([]*models.Template, error) {
	return m.dao.ListByGroupID(ctx, groupID)
}

func (m *manager) DeleteByID(ctx context.Context, id uint) error {
	return m.dao.DeleteByID(ctx, id)
}

func (m *manager) GetByID(ctx context.Context, id uint) (*models.Template, error) {
	return m.dao.GetByID(ctx, id)
}

func (m *manager) GetByName(ctx context.Context, name string) (*models.Template, error) {
	return m.dao.GetByName(ctx, name)
}

func (m *manager) GetRefOfApplication(ctx context.Context, id uint) ([]*amodels.Application, uint, error) {
	return m.dao.GetRefOfApplication(ctx, id)
}

func (m *manager) GetRefOfCluster(ctx context.Context, id uint) ([]*cmodels.Cluster, uint, error) {
	return m.dao.GetRefOfCluster(ctx, id)
}

func (m *manager) UpdateByID(ctx context.Context, id uint, template *models.Template) error {
	return m.dao.UpdateByID(ctx, id, template)
}

func (m *manager) ListByGroupIDs(ctx context.Context, ids []uint) ([]*models.Template, error) {
	return m.dao.ListByGroupIDs(ctx, ids)
}

func (m *manager) ListByIDs(ctx context.Context, ids []uint) ([]*models.Template, error) {
	return m.dao.ListByIDs(ctx, ids)
}

func (m *manager) ListV2(ctx context.Context, query *q.Query, groupIDs ...uint) ([]*models.Template, error) {
	return m.dao.ListV2(ctx, query, groupIDs...)
}
