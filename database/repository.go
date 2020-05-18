package database

func (db *DB) GetOffersList(queryString string , queryArgs []interface{}) ([]Offer, error) {
	offersList := make([]Offer, 0, 10)
	offer := Offer{}
	query := "SELECT short_id, title, content, created_at from offer "

	rows, err := db.Db.Query(query + queryString, queryArgs...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&offer.Id, &offer.Title, &offer.Content, &offer.CreatedAt)
		if err != nil {
			return nil, err
		}
		offersList = append(offersList, offer)
	}

	return offersList, err
}

func (db *DB) GetOffer(id string) (*Offer, error) {
	offer := &Offer{}
	query := "SELECT short_id, title, content, created_at from offer where short_id = $1"
	err := db.Db.QueryRow(query, id).Scan(&offer.Id, &offer.Title, &offer.Content, &offer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return offer, err
}

func (db *DB) GetCompany(id string) (*Company, error) {
	company := Company{}
	query := "SELECT c.short_id, COALESCE(ci.name, ''), COALESCE(ci.description, '') from company c left join company_info ci on c.info_id = ci.id where c.short_id = $1 "
	err := db.Db.QueryRow(query, id).Scan(&company.Id, &company.Name, &company.Description)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (db *DB) GetCompaniesList(queryString string , queryArgs []interface{}) ([]Company, error) {
	companiesList := make([]Company, 0, 10)
	company := Company{}
	query := "SELECT c.short_id, ci.name, ci.description from company c left join company_info ci on c.info_id = ci.id"

	rows, err := db.Db.Query(query + queryString, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&company.Id, &company.Name, &company.Description)
		if err != nil {
			return nil, err
		}
		companiesList = append(companiesList, company)
	}
	return companiesList, nil
}
