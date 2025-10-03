package repository

import (
	"rip/internal/app/ds"

	log "github.com/sirupsen/logrus"
)

func (r *TurbinesRepository) GetTurbines() (*[]ds.Turbine, error) {
	turbines := &[]ds.Turbine{}

	if err := r.db.Model(&ds.Turbine{}).Find(turbines).Error; err != nil {
		log.WithError(err).Error("Failed to get turbines")
		return nil, err
	}

	return turbines, nil
}

func (r *TurbinesRepository) GetTurbine(id uint) (*ds.Turbine, error) {
	turbine := &ds.Turbine{}

	if err := r.db.Where("id = ?", id).First(turbine).Error; err != nil {
		log.WithError(err).Error("Failed to get turbine")
		return nil, err
	}

	return turbine, nil
}

func (r *TurbinesRepository) GetTurbinesByTitle(title string) (*[]ds.Turbine, error) {
	turbines := &[]ds.Turbine{}

	if err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(turbines).Error; err != nil {
		log.WithError(err).Error("Failed to get turbines by title")
		return nil, err
	}

	return turbines, nil
}
