package contractions

type Repository interface {
	GetByID(id int) (interface{}, error)
	Update(model *interface{}) error
	Delete(model *interface{}) error
}
