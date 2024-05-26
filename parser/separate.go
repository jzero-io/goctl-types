package parser

func separateCommonElements(arrays ...[]string) [][]string {
	if len(arrays) == 0 {
		return nil
	}

	elementCount := make(map[string]int)

	for _, arr := range arrays {
		for _, elem := range arr {
			elementCount[elem]++
		}
	}

	uniqueArrays := make([][]string, len(arrays))

	for i, arr := range arrays {
		uniqueArr := make([]string, 0)
		for _, elem := range arr {
			if elementCount[elem] <= 1 {
				uniqueArr = append(uniqueArr, elem)
			}
		}
		uniqueArrays[i] = uniqueArr
	}

	return uniqueArrays
}
