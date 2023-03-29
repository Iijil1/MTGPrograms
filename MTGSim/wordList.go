package main

func listContains(list []string, word string) bool {
	for _, content := range list {
		if content == word {
			return true
		}
	}
	return false
}

func listsIntersect(list1, list2 []string) bool {
	for _, word1 := range list1 {
		if listContains(list2, word1) {
			return true
		}
	}
	return false
}

func removeFromList[T any](list *[]T, index int) {
	end := len(*list) - 1
	(*list)[index] = (*list)[end]
	*list = (*list)[:end]
}
