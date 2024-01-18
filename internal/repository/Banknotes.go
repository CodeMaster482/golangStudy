package repository

import (
	"main/internal/ds"
)

func (r *Repository) BanknoteList() (*[]ds.Banknote, error) {
	var banknotes []ds.Banknote
	if err := r.db.Where("is_delete = ?", false).Find(&banknotes).Error; err != nil {
		return nil, err
	}
	return &banknotes, nil
}

/*
func (r *Repository) SearchBanknotes(search string) (*[]ds.Banknote, error) {
	var banknotes []ds.Banknote
	r.db.Find(&banknotes)

	var filteredBanknotes []ds.Banknote
	for _, banknote := range banknotes {
		if strings.Contains(strings.ToLower(string(banknote.Amount)), strings.ToLower(search)) {
			filteredBanknotes = append(filteredBanknotes, banknote)
		}
	}

	return &filteredBanknotes, nil
}
*/

func (r *Repository) BanknoteById(id int) (*ds.Banknote, error) {
	var banknote ds.Banknote
	if err := r.db.First(&banknote, id).Error; err != nil {
		return nil, err
	}
	return &banknote, nil
}

func (r *Repository) DeleteBanknote(id string) error {
	query := "UPDATE banknotes SET is_delete=t WHERE id=$1;"
	result := r.db.Exec(query, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
