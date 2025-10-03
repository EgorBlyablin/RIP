package repository

import (
	"errors"
	"rip/internal/app/ds"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *TurbinesRepository) GetGenerationRequest(userId uint) (*ds.GenerationRequest, error) {
	generationRequest := &ds.GenerationRequest{}

	if err := r.db.Where(ds.GenerationRequest{
		CreatedByID: userId,
		Status:      "draft",
	}).Preload("TurbineGenerationRequests").Preload("TurbineGenerationRequests.Turbine").Attrs(ds.GenerationRequest{
		PeriodDays: 30,
	}).FirstOrCreate(generationRequest).Error; err != nil {
		log.WithError(err).Error("Failed to get or create generation request")
		return nil, err
	}

	return generationRequest, nil
}

func (r *TurbinesRepository) GetGenerationRequestsCount(userId uint) uint {
	generationRequests := &[]ds.GenerationRequest{}
	err := r.db.
		Where("created_by_id = ? AND status = ?", userId, "draft").
		Order("id DESC").
		First(&generationRequests).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0
	}
	if err != nil {
		log.WithError(err).Error("Failed to get generation request")
		return 0
	}

	return uint(r.db.Model(*generationRequests).Association("TurbineGenerationRequests").Count())
}

func (r *TurbinesRepository) AddTurbineToGenerationRequest(userId uint, turbineId uint) error {
	generationRequest, err := r.GetGenerationRequest(userId)
	if err != nil {
		log.WithError(err).Error("Failed to get generation request")
		return err
	}

	r.db.Where(&ds.TurbineGenerationRequest{TurbineID: turbineId, GenerationRequestID: generationRequest.ID}).FirstOrCreate(&ds.TurbineGenerationRequest{})

	return nil
}

func (r *TurbinesRepository) DeleteGenerationRequest(userId uint) error {
	if err := r.db.Exec(
		"UPDATE generation_requests SET status = ?, deleted_at = ? WHERE created_by_id = ? AND status = ?",
		"deleted",
		time.Now(),
		userId,
		"draft",
	).Error; err != nil {
		log.WithError(err).Error("Failed to delete generation request")
		return err
	}
	return nil
}
