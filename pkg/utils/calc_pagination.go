package utils

func SafeCalculate(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	if pageSize > 100 {
		pageSize = 100
	}

	limit = pageSize
	offset = (page - 1) * pageSize

	return offset, limit
}

func SafeCalculateForPostgres(page, pageSize int) (offset, limit int) {
	offset, limit = SafeCalculate(page, pageSize)

	if offset > 10000 {
		offset = 10000
	}

	return offset, limit
}
