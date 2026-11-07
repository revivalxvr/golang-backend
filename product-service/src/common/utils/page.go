package utils

func CreateLimitAndOffset(page int) (limit, offset int) {
	const limitDefault = 50

	if page <= 0 {
		page = 1
	}
	limit = limitDefault
	offset = (page - 1) * limit
	return
}
