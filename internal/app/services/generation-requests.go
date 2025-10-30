package services

import (
	"errors"
	"math"

	"rip/internal/app/ds"
	"rip/internal/app/repositories"

	"gorm.io/gorm"
)

var (
	ErrorGenerationRequestIsNotDraft            = errors.New("generation request is not in \"draft\" status")
	ErrorGenerationRequestPeriodDaysRequired    = errors.New("period_days must be filled before generation request can be processed")
	ErrorGenerationRequestIncompleteTurbineData = errors.New("turbine data is incomplete for generation calculation")
	ErrorGenerationRequestDraftHasNoTurbines    = errors.New("generation request draft must contain at least one turbine")
	ErrorGenerationRequestCannotBeDeleted       = errors.New("only draft generation requests can be deleted")
	ErrorGenerationRequestIncorrectStatus       = errors.New("passed generation request status is incorrect")
)

type GenerationRequestsService struct {
	GenerationRequestsRepository *repositories.GenerationRequestRepository
	usersService                 *UsersService
}

func CalculateTurbineGeneration(avgVelocity float32, height uint16, alpha float32, power uint32, days uint) uint64 {
	const (
		optimalVelocity     = 12
		cutoffVelocity      = 15
		decreaseCoefficient = 0.4
	)

	velocityAtHeight := float64(avgVelocity) * math.Pow(float64(height)/10.0, float64(alpha))

	generation := float64(power) *
		math.Pow(velocityAtHeight/float64(optimalVelocity), 3) *
		math.Exp(-decreaseCoefficient*velocityAtHeight/float64(cutoffVelocity)) *
		24 * float64(days)

	return uint64(generation)
}

func NewGenerationRequestsService(turbinesAppDB *gorm.DB) *GenerationRequestsService {
	return &GenerationRequestsService{
		GenerationRequestsRepository: repositories.NewGenerationRequestRepository(turbinesAppDB),
		usersService:                 NewUsersService(turbinesAppDB),
	}
}

func (service *GenerationRequestsService) GetGenerationRequests(userId uint, filter repositories.GenerationRequestsFilter) ([]ds.GenerationRequest, error) {
	return service.GenerationRequestsRepository.GetGenerationRequests(userId, filter)
}

func (service *GenerationRequestsService) GetGenerationRequest(generationRequestId uint) (ds.GenerationRequest, error) {
	return service.GenerationRequestsRepository.GetGenerationRequest(generationRequestId)
}

func (service *GenerationRequestsService) CloseGenerationRequest(generationRequestId uint, userId uint, status string) (ds.GenerationRequest, error) {
	usersService := NewUsersService(service.GenerationRequestsRepository.GenerationRequestDB)
	isModerator, err := usersService.CheckIsModerator(userId)
	if err != nil {
		return ds.GenerationRequest{}, err
	} else if !isModerator {
		return ds.GenerationRequest{}, ErrorUserHasNoAccess
	}

	switch status {
	case "completed":
		return service.GenerationRequestsRepository.CompleteGenerationRequest(generationRequestId, userId, CalculateTurbineGeneration)
	case "rejected":
		return service.GenerationRequestsRepository.RejectGenerationRequest(generationRequestId, userId)
	}

	return ds.GenerationRequest{}, ErrorGenerationRequestIncorrectStatus
}

func (service *GenerationRequestsService) GetDraftBriefInfo(userId uint) (ds.DraftGenerationRequestsBriefInfo, error) {
	currentDraft, err := service.GenerationRequestsRepository.GetDraftGenerationRequest(userId)
	if err != nil {
		return ds.DraftGenerationRequestsBriefInfo{}, err
	}

	currentDraftTurbinesCount, err := service.GenerationRequestsRepository.GetGenerationRequestsTurbinesCount(currentDraft.ID)
	if err != nil {
		return ds.DraftGenerationRequestsBriefInfo{}, err
	}

	return ds.DraftGenerationRequestsBriefInfo{
		GenerationRequestId: currentDraft.ID,
		TurbinesCount:       currentDraftTurbinesCount,
	}, nil
}

func (service *GenerationRequestsService) UpdateDraftGenerationRequest(userId uint, generationRequestUpdates ds.UpdateGenerationRequest) (ds.GenerationRequest, error) {
	return service.GenerationRequestsRepository.UpdateDraftGenerationRequest(userId, generationRequestUpdates)
}

func (service *GenerationRequestsService) AddTurbineToDraft(userId, turbineId uint) error {
	draftGenerationRequest, err := service.GenerationRequestsRepository.GetOrCreateDraftGenerationRequest(userId)
	if err != nil {
		return err
	}

	return service.GenerationRequestsRepository.AddTurbineToDraftGenerationRequest(draftGenerationRequest.ID, turbineId)
}

func (service *GenerationRequestsService) UpdateTurbineInDraft(userId, turbineId uint, updates ds.UpdateTurbineGenerationRequest) error {
	currentDraft, err := service.GenerationRequestsRepository.GetDraftGenerationRequest(userId)
	if err != nil {
		return err
	}

	return service.GenerationRequestsRepository.UpdateTurbineInDraft(currentDraft.ID, turbineId, updates)
}

func (service *GenerationRequestsService) RemoveTurbineFromDraft(userId, turbineId uint) error {
	currentDraft, err := service.GenerationRequestsRepository.GetDraftGenerationRequest(userId)
	if err != nil {
		return err
	}

	return service.GenerationRequestsRepository.RemoveTurbineFromDraftGenerationRequest(currentDraft.ID, turbineId)
}

func (service *GenerationRequestsService) SubmitDraftGenerationRequest(userId uint) error {
	return service.GenerationRequestsRepository.SubmitDraftGenerationRequest(userId)
}

func (service *GenerationRequestsService) DeleteDraftGenerationRequest(userId uint) error {
	currentDraft, err := service.GenerationRequestsRepository.GetDraftGenerationRequest(userId)
	if err != nil {
		return err
	}

	return service.GenerationRequestsRepository.DeleteGenerationRequest(currentDraft.ID)
}
