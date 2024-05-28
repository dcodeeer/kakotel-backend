package estates

import (
	"api/internal/core"
	"api/internal/infrastructure"
)

type estates struct {
	repo infrastructure.IEstates
}

func New(repo infrastructure.IEstates) *estates {
	return &estates{repo: repo}
}

func (s *estates) Add(estate *core.Estate) (int, error) {
	id, err := s.repo.Add(estate)

	s.repo.AddAddress(&core.Address{
		EstateId: id,
		Number:   estate.Address.Number,
		Street:   estate.Address.Street,
		City:     estate.Address.City,
		District: estate.Address.District,
	})

	return id, err
}

func (s *estates) AddTempImage(bytes []byte) (int, error) {
	filename, err := core.CreateFile(bytes)
	if err != nil {
		return 0, err
	}
	return s.repo.AddTempImage(filename)
}

func (s *estates) GetTempImages(ids []int) ([]string, error) {
	return s.repo.GetTempImages(ids)
}
func (s *estates) GetCategories() ([]core.Category, error) {
	return s.repo.GetCategories()
}

func (s *estates) GetAmenities() ([]core.Amenity, error) {
	return s.repo.GetAmenities()
}

func (s *estates) GetAll() (*[]core.Estate, error) {
	return s.repo.GetAll()
}

func (s *estates) GetById(id int) (*core.Estate, error) {
	return s.repo.GetOne(id)
}
