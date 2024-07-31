package services

import (
	"Messaggio/producer/internal/model"
	"fmt"
)

func (s *Services) SaveMessage(message *model.Message) (string, error) {
	return s.repository.SaveMessage(message)
}

func (s *Services) GetMessage(id string) (*model.Message, error) {
	return s.repository.GetMessage(id)
}

func (s *Services) GetMessages() ([]*model.Message, error) {
	return s.repository.GetMessages()
}

func (s *Services) DeleteMessages() error {
	return s.repository.DeleteMessages()
}

func (s *Services) DeleteMessage(id string) error {
	return s.repository.DeleteMessage(id)
}

func (s *Services) GetStatsByDays() (map[string]*model.ByDays, error) {
	stats, err := s.repository.GetStatsByHour()
	if err != nil {
		return nil, err
	}

	ans := make(map[string]*model.ByDays)

	for _, stat := range stats {
		year, month, day := stat.HourStat.Date()
		key := fmt.Sprintf("%d-%d-%d", year, month, day)

		if _, ok := ans[key]; !ok {
			ans[key] = &model.ByDays{
				HourStats:    make(map[int]int),
				OverallCount: 0,
			}
		}

		ans[key].HourStats[stat.HourStat.Hour()] += stat.Count
		ans[key].OverallCount += stat.Count
	}

	return ans, nil
}
