package cluster

import (
	"context"
	"fmt"
	"time"

	"g.hz.netease.com/horizon/core/common"
	"g.hz.netease.com/horizon/core/controller/build"
	controllertag "g.hz.netease.com/horizon/core/controller/tag"
	herrors "g.hz.netease.com/horizon/core/errors"
	appgitrepo "g.hz.netease.com/horizon/pkg/application/gitrepo"
	appmodels "g.hz.netease.com/horizon/pkg/application/models"
	clustercommon "g.hz.netease.com/horizon/pkg/cluster/common"
	"g.hz.netease.com/horizon/pkg/cluster/gitrepo"
	"g.hz.netease.com/horizon/pkg/templaterelease/models"
	templateschema "g.hz.netease.com/horizon/pkg/templaterelease/schema"
	"g.hz.netease.com/horizon/pkg/util/jsonschema"
	"g.hz.netease.com/horizon/pkg/util/mergemap"

	codemodels "g.hz.netease.com/horizon/pkg/cluster/code"
	perror "g.hz.netease.com/horizon/pkg/errors"
	regionmodels "g.hz.netease.com/horizon/pkg/region/models"

	"g.hz.netease.com/horizon/pkg/util/wlog"
)

func (c *controller) CreateClusterV2(ctx context.Context, applicationID uint, environment,
	region string, r *CreateClusterRequestV2, mergePatch bool) (*CreateClusterResponseV2, error) {
	const op = "cluster controller: create cluster v2"
	defer wlog.Start(ctx, op).StopPrint()

	// 1. check exist
	exists, err := c.clusterMgr.CheckClusterExists(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, perror.Wrap(herrors.ErrNameConflict,
			"a cluster with the same name already exists, please do not create it again")
	}

	// 2. validate create req
	if err := validateClusterName(r.Name); err != nil {
		return nil, err
	}

	// 3. get application
	application, err := c.applicationMgr.GetByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	// 4. customize buildTemplateInfo and do validate
	buildTemplateInfo, err := c.customizeCreateReqBuildTemplateInfo(ctx, r, application, environment, mergePatch)
	if err != nil {
		return nil, err
	}

	if err := buildTemplateInfo.Validate(ctx,
		c.templateSchemaGetter, nil, c.buildSchema); err != nil {
		return nil, err
	}

	// 5. get environment and region
	envEntity, err := c.envRegionMgr.GetByEnvironmentAndRegion(ctx, environment, region)
	if err != nil {
		return nil, err
	}
	regionEntity, err := c.regionMgr.GetRegionEntity(ctx, region)
	if err != nil {
		return nil, err
	}
	if regionEntity.Disabled {
		return nil, perror.Wrap(herrors.ErrDisabled,
			"the region which is disabled cannot be used to create a cluster")
	}

	// 6. transfer expireTime to expireSeconds and verify environment.
	expireSeconds, err := c.toExpireSeconds(ctx, r.ExpireTime, environment)
	if err != nil {
		return nil, err
	}

	// 7. get templateRelease
	tr, err := c.templateReleaseMgr.GetByTemplateNameAndRelease(ctx, r.TemplateInfo.Name, r.TemplateInfo.Release)
	if err != nil {
		return nil, err
	}

	// 8. customize db infos
	cluster, tags := r.toClusterModel(application, envEntity, buildTemplateInfo, expireSeconds)

	// 9. update db and tags
	clusterResp, err := c.clusterMgr.Create(ctx, cluster, tags, r.ExtraMembers)
	if err != nil {
		return nil, err
	}

	// 10. create git repo
	err = c.clusterGitRepo.CreateCluster(ctx, &gitrepo.CreateClusterParams{
		BaseParams: &gitrepo.BaseParams{
			ClusterID:           clusterResp.ID,
			Cluster:             clusterResp.Name,
			PipelineJSONBlob:    r.BuildConfig,
			ApplicationJSONBlob: r.TemplateConfig,
			TemplateRelease:     tr,
			Application:         application,
			Environment:         environment,
			RegionEntity:        regionEntity,
			Version:             common.MetaVersion2,
		},
		Tags: tags,
	})
	if err != nil {
		// Prevent errors like "project has already been taken" caused by automatic retries due to api timeouts
		if deleteErr := c.clusterGitRepo.DeleteCluster(ctx, application.Name, cluster.Name, cluster.ID); deleteErr != nil {
			if _, ok := perror.Cause(deleteErr).(*herrors.HorizonErrNotFound); !ok {
				err = perror.WithMessage(err, deleteErr.Error())
			}
		}
		if deleteErr := c.clusterMgr.DeleteByID(ctx, cluster.ID); deleteErr != nil {
			err = perror.WithMessage(err, deleteErr.Error())
		}
		return nil, err
	}
	cluster.Status = clustercommon.StatusEmpty
	updateClusterResp, err := c.clusterMgr.UpdateByID(ctx, cluster.ID, cluster)
	if err != nil {
		return nil, err
	}

	// 11. get full path
	group, err := c.groupSvc.GetChildByID(ctx, application.GroupID)
	if err != nil {
		return nil, err
	}
	fullPath := fmt.Sprintf("%v/%v/%v", group.FullPath, application.Name, cluster.Name)

	// 12. customize response
	return &CreateClusterResponseV2{
		ID:            updateClusterResp.ID,
		FullPath:      fullPath,
		ApplicationID: applicationID,
		CreatedAt:     updateClusterResp.CreatedAt,
		UpdatedAt:     updateClusterResp.UpdatedAt,
	}, nil
}

func (c *controller) GetClusterV2(ctx context.Context, clusterID uint) (*GetClusterResponseV2, error) {
	const op = "cluster controller: get cluster v2"
	defer wlog.Start(ctx, op).StopPrint()

	// 1. get cluster from db
	cluster, err := c.clusterMgr.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	// 2. get application
	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return nil, err
	}

	// 3. get region entity
	regionEntity, err := c.regionMgr.GetRegionEntity(ctx, cluster.RegionName)
	if err != nil {
		return nil, err
	}

	// 4. gen fullPath
	fullPath, err := func() (string, error) {
		group, err := c.groupSvc.GetChildByID(ctx, application.GroupID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v/%v/%v", group.FullPath, application.Name, cluster.Name), nil
	}()
	if err != nil {
		return nil, err
	}

	// 5. get tags
	tags, err := c.tagMgr.ListByResourceTypeID(ctx, common.ResourceCluster, cluster.ID)
	if err != nil {
		return nil, err
	}

	// 6. get GitRepo
	clusterGitRepoFile, err := c.clusterGitRepo.GetCluster(ctx, application.Name, cluster.Name, cluster.Template)
	if err != nil {
		return nil, err
	}

	// 7. get createdBy and updatedBy users
	userMap, err := c.userManager.GetUserMapByIDs(ctx, []uint{cluster.CreatedBy, cluster.UpdatedBy})
	if err != nil {
		return nil, err
	}
	getResp := &GetClusterResponseV2{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Description: cluster.Description,
		// TODO: currently cluster not support different priority with application
		Priority: string(application.Priority),
		Scope: &Scope{
			Environment:       cluster.EnvironmentName,
			Region:            cluster.RegionName,
			RegionDisplayName: regionEntity.DisplayName,
		},
		ExpireTime: func() string {
			expireTime := ""
			if cluster.ExpireSeconds > 0 {
				expireTime = time.Duration(cluster.ExpireSeconds * 1e9).String()
			}
			return expireTime
		}(),
		FullPath:        fullPath,
		ApplicationName: application.Name,
		ApplicationID:   application.ID,
		Tags: func() []controllertag.Tag {
			retTags := make([]controllertag.Tag, 0, len(tags))
			for _, tag := range retTags {
				retTags = append(retTags, controllertag.Tag{
					Key:   tag.Key,
					Value: tag.Value,
				})
			}
			return retTags
		}(),
		Git: func() *codemodels.Git {
			if cluster.GitURL == "" {
				return nil
			}
			return codemodels.NewGit(cluster.GitURL, cluster.GitSubfolder,
				cluster.GitRefType, cluster.GitRef)
		}(),
		BuildConfig: clusterGitRepoFile.PipelineJSONBlob,
		TemplateInfo: func() *codemodels.TemplateInfo {
			if cluster.Template == "" {
				return nil
			}
			return &codemodels.TemplateInfo{
				Name:    cluster.Template,
				Release: cluster.TemplateRelease,
			}
		}(),
		TemplateConfig: clusterGitRepoFile.ApplicationJSONBlob,
		Manifest:       clusterGitRepoFile.Manifest,
		Status:         cluster.Status,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
		CreatedBy:      toUser(getUserFromMap(cluster.CreatedBy, userMap)),
		UpdatedBy:      toUser(getUserFromMap(cluster.UpdatedBy, userMap)),
	}
	return getResp, nil
}

func (c *controller) UpdateClusterV2(ctx context.Context, clusterID uint,
	r *UpdateClusterRequestV2, mergePatch bool) error {
	const op = "cluster controller: update cluster v2"
	defer wlog.Start(ctx, op).StopPrint()

	// 1. get cluster and application from db
	cluster, err := c.clusterMgr.GetByID(ctx, clusterID)
	if err != nil {
		return err
	}
	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return err
	}

	// 2. check if we should update region and env
	var regionEntity *regionmodels.RegionEntity
	environmentName := cluster.EnvironmentName
	regionName := cluster.RegionName
	if cluster.Status == clustercommon.StatusFreed &&
		(r.Environment != nil && *r.Environment != environmentName) ||
		(r.Region != nil && *r.Region != regionName) {
		if r.Environment != nil {
			environmentName = *r.Environment
		}
		if r.Region != nil {
			regionName = *r.Region
		}
		regionEntity, err = c.regionMgr.GetRegionEntity(ctx, regionName)
		if err != nil {
			return err
		}
		_, err = c.envRegionMgr.GetByEnvironmentAndRegion(ctx, environmentName, regionName)
		if err != nil {
			return err
		}
	}

	// 3. check and transfer ExpireTime
	expireSeconds := cluster.ExpireSeconds
	if r.ExpireTime != "" {
		expireSeconds, err = c.toExpireSeconds(ctx, r.ExpireTime, environmentName)
		if err != nil {
			return err
		}
	}

	// 4. customize template\build\template infos
	templateInfo, templateRelease, err := func() (*codemodels.TemplateInfo, *models.TemplateRelease, error) {
		var templateInfo *codemodels.TemplateInfo
		if r.TemplateInfo == nil {
			templateInfo = &codemodels.TemplateInfo{
				Name:    cluster.Template,
				Release: cluster.TemplateRelease,
			}
		} else {
			templateInfo = r.TemplateInfo
		}
		tr, err := c.templateReleaseMgr.GetByTemplateNameAndRelease(ctx,
			templateInfo.Name, templateInfo.Release)
		if err != nil {
			return nil, nil, err
		}
		return templateInfo, tr, nil
	}()
	if err != nil {
		return err
	}

	buildConfig, templateConfig, err := func() (map[string]interface{}, map[string]interface{}, error) {
		if r.BuildConfig == nil && r.TemplateConfig == nil {
			return nil, nil, nil
		}
		files, err := c.clusterGitRepo.GetCluster(ctx, application.Name, cluster.Name, cluster.Template)
		if err != nil {
			return nil, nil, err
		}
		if files.Manifest == nil {
			return nil, nil, perror.Wrapf(herrors.ErrParamInvalid, "git repo  %s not support v2 interface",
				cluster.Name)
		}

		buildConfig := r.BuildConfig
		templateConfig := r.TemplateConfig
		if r.BuildConfig != nil && mergePatch {
			buildConfig, err = mergemap.Merge(files.PipelineJSONBlob, r.BuildConfig)
			if err != nil {
				return nil, nil, err
			}
		}
		if r.TemplateConfig != nil && mergePatch {
			templateConfig, err = mergemap.Merge(files.ApplicationJSONBlob, r.TemplateConfig)
			if err != nil {
				return nil, nil, err
			}
		}
		return buildConfig, templateConfig, nil
	}()
	if err != nil {
		return err
	}

	// 5. validate update Request
	err = func() error {
		renderValues, err := c.getRenderValueFromTag(ctx, clusterID)
		if err != nil {
			return err
		}
		info := BuildTemplateInfo{
			BuildConfig:    buildConfig,
			TemplateInfo:   templateInfo,
			TemplateConfig: templateConfig,
		}
		return info.Validate(ctx, c.templateSchemaGetter, renderValues, c.buildSchema)
	}()
	if err != nil {
		return nil
	}

	// 6. update in git repo
	if err = c.clusterGitRepo.UpdateCluster(ctx, &gitrepo.UpdateClusterParams{
		BaseParams: &gitrepo.BaseParams{
			ClusterID:           cluster.ID,
			Cluster:             cluster.Name,
			PipelineJSONBlob:    buildConfig,
			ApplicationJSONBlob: templateConfig,
			TemplateRelease:     templateRelease,
			Application:         application,
			Environment:         environmentName,
			RegionEntity:        regionEntity,
			Version:             common.MetaVersion2,
		}}); err != nil {
		return err
	}

	// 7. update cluster in db
	clusterModel := r.toClusterModel(cluster, expireSeconds, environmentName,
		regionName, templateInfo.Name, templateInfo.Release)
	_, err = c.clusterMgr.UpdateByID(ctx, clusterID, clusterModel)
	if err != nil {
		return err
	}
	return nil
}

type BuildTemplateInfo struct {
	BuildConfig    map[string]interface{}
	TemplateInfo   *codemodels.TemplateInfo
	TemplateConfig map[string]interface{}
}

func (info *BuildTemplateInfo) Validate(ctx context.Context,
	trGetter templateschema.Getter, templateSchemaRenderVal map[string]string, buildSchema *build.Schema) error {
	if templateSchemaRenderVal == nil {
		templateSchemaRenderVal = make(map[string]string)
	}
	// TODO (remove it, currently some template need it)
	templateSchemaRenderVal["resourceType"] = "cluster"
	schema, err := trGetter.GetTemplateSchema(ctx, info.TemplateInfo.Name,
		info.TemplateInfo.Release, templateSchemaRenderVal)
	if err != nil {
		return err
	}
	err = jsonschema.Validate(schema.Application.JSONSchema, info.TemplateConfig, false)
	if err != nil {
		return err
	}

	if buildSchema != nil && info.BuildConfig != nil {
		err = jsonschema.Validate(buildSchema.JSONSchema, info.BuildConfig, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) customizeCreateReqBuildTemplateInfo(ctx context.Context, r *CreateClusterRequestV2,
	application *appmodels.Application, environment string, mergePatch bool) (*BuildTemplateInfo, error) {
	buildTemplateInfo := &BuildTemplateInfo{}

	var appGitRepoFile *appgitrepo.GetResponse
	var err error
	if (r.BuildConfig != nil || r.TemplateInfo != nil) && mergePatch {
		appGitRepoFile, err = c.applicationGitRepo.GetApplication(ctx, application.Name, environment)
		if err != nil {
			return nil, err
		}
	}
	if r.BuildConfig != nil {
		if mergePatch {
			buildTemplateInfo.BuildConfig, err = mergemap.Merge(appGitRepoFile.BuildConf, r.BuildConfig)
			if err != nil {
				return nil, err
			}
		} else {
			buildTemplateInfo.BuildConfig = r.BuildConfig
		}
	} else {
		buildTemplateInfo.BuildConfig = appGitRepoFile.BuildConf
	}

	if r.TemplateInfo == nil && r.TemplateConfig == nil {
		buildTemplateInfo.TemplateInfo = &codemodels.TemplateInfo{
			Name:    application.Template,
			Release: application.TemplateRelease,
		}
		buildTemplateInfo.TemplateConfig = appGitRepoFile.TemplateConf
	} else if r.TemplateInfo != nil && r.TemplateConfig != nil {
		buildTemplateInfo.TemplateInfo = r.TemplateInfo
		if mergePatch {
			buildTemplateInfo.TemplateConfig, err = mergemap.Merge(appGitRepoFile.TemplateConf, r.TemplateConfig)
			if err != nil {
				return nil, err
			}
		} else {
			buildTemplateInfo.TemplateConfig = r.TemplateConfig
		}
	} else {
		return nil, perror.Wrap(herrors.ErrParamInvalid, "TemplateInfo or TemplateConfig nil")
	}
	return buildTemplateInfo, nil
}