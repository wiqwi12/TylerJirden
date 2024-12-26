package application

import "GymBot/internal/domain/repository"

type Service struct {
	Repo repository.UserRepository
	//Analitycs repository.AnalitycsRepository
}

func Initialize(repo repository.UserRepository) *Service {
	return &Service{
		Repo: repo,
		//	Analitycs: analitycs,
	}
}
