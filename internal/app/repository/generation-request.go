package repository

import (
	"errors"
	"rip/internal/app/ds"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *TurbinesRepository) GetGenerationRequest(generationRequestId uint) (*ds.GenerationRequest, error) {
	generationRequest := &ds.GenerationRequest{}

	err := r.db.Where(&ds.GenerationRequest{ID: generationRequestId}).Preload("TurbineGenerationRequests").Preload("TurbineGenerationRequests.Turbine").First(generationRequest).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return generationRequest, nil
}

func (r *TurbinesRepository) GetDraftGenerationRequest(userId uint) (*ds.GenerationRequest, error) {
	generationRequest := &ds.GenerationRequest{}

	err := r.db.Where(&ds.GenerationRequest{
		CreatedByID: userId,
		Status:      "draft",
	}).Preload("TurbineGenerationRequests").Preload("TurbineGenerationRequests.Turbine").First(generationRequest).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return generationRequest, nil
}

func (r *TurbinesRepository) GetOrCreateUserDraftGenerationRequest(userId uint) (*ds.GenerationRequest, error) {
	generationRequest := &ds.GenerationRequest{}

	if err := r.db.Where(&ds.GenerationRequest{
		CreatedByID: userId,
		Status:      "draft",
	}).Preload("TurbineGenerationRequests").Preload("TurbineGenerationRequests.Turbine").FirstOrCreate(generationRequest).Error; err != nil {
		log.WithError(err).Error("Failed to get or create draft generation request")
		return nil, err
	}

	return generationRequest, nil
}

func (r *TurbinesRepository) GetDraftGenerationRequestsCount(userId uint) uint {
	generationRequests := &[]ds.GenerationRequest{}
	err := r.db.Where("created_by_id = ? AND status = ?", userId, "draft").First(&generationRequests).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0
	}
	if err != nil {
		log.WithError(err).Error("Failed to get generation request")
		return 0
	}

	return uint(r.db.Model(*generationRequests).Association("TurbineGenerationRequests").Count())
}

func (r *TurbinesRepository) AddTurbineToDraftGenerationRequest(generationRequestId uint, turbineId uint) error {
	r.db.Where(&ds.TurbineGenerationRequest{
		TurbineID:           turbineId,
		GenerationRequestID: generationRequestId,
	}).FirstOrCreate(&ds.TurbineGenerationRequest{})

	return nil
}

func (r *TurbinesRepository) DeleteDraftGenerationRequest(generationRequestId uint) error {
	if err := r.db.Exec(
		"UPDATE generation_requests SET status = ?, deleted_at = ? WHERE id = ? AND status = ?",
		"deleted",
		time.Now(),
		generationRequestId,
		"draft",
	).Error; err != nil {
		log.WithError(err).Error("Failed to delete generation request")
		return err
	}
	return nil
}
