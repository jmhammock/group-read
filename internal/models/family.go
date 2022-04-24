package models

type Family struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Families []Family

type FamilyModel struct {
	db DB
}

func NewFamilyModel(db DB) *FamilyModel {
	return &FamilyModel{
		db: db,
	}
}

func (f FamilyModel) Get(limit, offset uint16) (*Families, error) {
	q := `SELECT *
		FROM families
		LIMIT $1
		OFFSET $2;`

	rows, err := f.db.Query(q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var families Families
	for rows.Next() {
		var family Family
		err := rows.Scan(&family)
		if err != nil {
			return nil, err
		}
		families = append(families, family)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &families, nil
}

func (f FamilyModel) GetById(id string) (*Family, error) {
	q := `SELECT *
		FROM families
		WHERE id = $1;`

	var family *Family
	err := f.db.QueryRow(q, id).Scan(&family)

	return family, err
}

func (f FamilyModel) Insert(name string) (string, error) {
	q := `INSERT INTO families (name) VALUES($1)
		RETURNING id;`

	var id string
	err := f.db.QueryRow(q, name).Scan(&id)

	return id, err
}

func (f FamilyModel) Update(name string) error {
	q := `UPDATE families
		SET name = $2
		WHERE id = $1;`

	_, err := f.db.Execute(q, name)

	return err
}

func (f FamilyModel) Delete(id string) error {
	q := `DELETE FROM families
		WHERE id = $1;`

	_, err := f.db.Execute(q, id)
	return err
}
