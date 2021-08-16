package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/repository"
	"strconv"
)

// Repository struct
type Repository Controller

// GetCommitList get latest 10 commit list
func (Repository) GetCommitList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo("git")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	list, err := repo.BranchLog(project.ID, gp.URLQuery.Get("branch"), 10)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			CommitList []repository.CommitInfo `json:"list"`
		}{CommitList: list},
	}
}

// GetBranchList get all branch list
func (Repository) GetBranchList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo("git")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	branchList, err := repo.BranchList(project.ID)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			BranchList []string `json:"list"`
		}{BranchList: branchList},
	}
}

// GetTagList get latest 10 tag list
func (Repository) GetTagList(gp *core.Goploy) *core.Response {
	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	project, err := model.Project{ID: id}.GetData()
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	repo, err := repository.GetRepo("git")
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	list, err := repo.TagLog(project.ID, 10)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	return &core.Response{
		Data: struct {
			TagList []repository.CommitInfo `json:"list"`
		}{TagList: list},
	}
}
