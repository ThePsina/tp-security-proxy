package application

import (
	"proxy/internal/domain/entity"
	"proxy/internal/domain/repository"
)

type DataManager struct {
	db repository.ReqRepo
}

func NewDataManager(db repository.ReqRepo) *DataManager {
	return &DataManager{db}
}

func (dm *DataManager) Insert(req entity.Req) error {
	return dm.db.Insert(req)
}

func (dm *DataManager) GetRequestList() ([]entity.Req, error) {
	return dm.db.GetRequestList()
}

func (dm *DataManager) GetRequestById(id int64) (entity.Req, error) {
	return dm.db.GetRequestById(id)
}

func (dm *DataManager) GetRequestHeaders(id int64) (entity.Req, error) {
	return dm.db.GetRequestHeaders(id)
}
