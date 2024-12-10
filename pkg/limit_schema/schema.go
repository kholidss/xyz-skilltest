package limit_schema

func BuildLimitSchema(salary int) map[int]int {
	if salary <= 5000000 {
		return map[int]int{
			1: 100000,
			2: 200000,
			3: 300000,
			6: 400000,
		}
	}

	if salary > 5000000 && salary <= 10000000 {
		return map[int]int{
			1: 500000,
			2: 600000,
			3: 700000,
			6: 800000,
		}
	}

	return map[int]int{
		1: 1000000,
		2: 1200000,
		3: 1500000,
		6: 2000000,
	}
}
