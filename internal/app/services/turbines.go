package services

import (
	"io"
	"rip/internal/app/ds"
	"rip/internal/app/repositories"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TurbinesService struct {
	turbinesRepository *repositories.TurbinesRepository
	turbinesImages     *repositories.S3Repository
}

func NewTurbinesService(turbinesAppDB *gorm.DB, turbinesAppS3 *repositories.S3Repository) *TurbinesService {
	return &TurbinesService{
		turbinesRepository: repositories.NewTurbinesRepository(turbinesAppDB),
		turbinesImages:     turbinesAppS3,
	}
}

func (service *TurbinesService) GetActiveTurbines(titleFilter *string) ([]ds.Turbine, error) {
	return service.turbinesRepository.GetActiveTurbines(repositories.TurbinesFilter{
		Title: titleFilter,
	})
}

func (service *TurbinesService) GetActiveTurbine(turbineId uint) (ds.Turbine, error) {
	return service.turbinesRepository.GetActiveTurbine(turbineId)
}

func (service *TurbinesService) CreateTurbine(turbine ds.CreateTurbine) (ds.Turbine, error) {
	return service.turbinesRepository.CreateTurbine(turbine)
}

func (service *TurbinesService) UpdateTurbine(turbineId uint, turbine ds.UpdateTurbine) (ds.Turbine, error) {
	_, err := service.GetActiveTurbine(turbineId)
	if err != nil {
		log.WithError(err).WithField("Turbine ID", turbineId).Error("Failed to get turbine during update process")
		return ds.Turbine{}, err
	}

	return service.turbinesRepository.UpdateTurbine(turbineId, turbine)
}

func (service *TurbinesService) UpdateTurbineImage(turbineId uint, image io.Reader, imageSize int64, imageContentType string) (ds.Turbine, error) {
	turbine, err := service.GetActiveTurbine(turbineId)
	if err != nil {
		log.WithError(err).WithField("Turbine ID", turbineId).Error("Failed to get turbine during image upload process")
		return ds.Turbine{}, err
	}

	turbineImageKey, err := service.turbinesImages.UploadFile(turbine.ImageSlug(), image, imageSize, imageContentType)
	if err != nil {
		log.WithError(err).WithField("Turbine ID", turbineId).Error("Failed to upload turbine image")
		return ds.Turbine{}, err
	}

	updatedTurbine, err := service.turbinesRepository.UpdateTurbine(turbineId, ds.UpdateTurbine{
		Image: &turbineImageKey,
	})
	return updatedTurbine, err
}

func (service *TurbinesService) DeleteTurbine(turbineId uint) error {
	turbine, err := service.GetActiveTurbine(turbineId)
	if err != nil {
		log.WithError(err).WithField("Turbine ID", turbineId).Error("Failed to get turbine during deletion process")
		return err
	}

	if turbine.Image != nil {
		err = service.turbinesImages.DeleteFile(turbine.ImageSlug())
		if err != nil {
			log.WithError(err).WithField("Turbine ID", turbineId).Error("Failed to delete turbine image")
			return err
		}
	}

	err = service.turbinesRepository.DeleteTurbine(turbineId)
	if err != nil {
		log.WithError(err).WithField("Turbine ID", turbineId).Error("Failed to delete turbine")
		return err
	}

	return nil
}
