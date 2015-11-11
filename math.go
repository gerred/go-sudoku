package main

func setsAreEqual(a []int, b []int) bool {
	store := make(map[int]interface{})
	for _, val := range a {
		store[val] = struct{}{}
	}
	for _, val := range b {
		if _, ok := store[val]; !ok {
			return false
		}
		delete(store, val)
	}
	return true
}

func intersect(a []int, b []int) []int {
	store := make(map[int]interface{})
	var list []int
	for _, i := range a {
		for _, j := range b {
			if i == j {
				if _, ok := store[i]; !ok {
					store[i] = struct{}{}
					list = append(list, i)
				}
			}
		}
	}
	return list
}

func union(a []int, b []int) []int {
	store := make(map[int]interface{})
	var list []int
	for _, val := range a {
		if _, ok := store[val]; !ok {
			store[val] = struct{}{}
			list = append(list, val)
		}
	}
	for _, val := range b {
		if _, ok := store[val]; !ok {
			store[val] = struct{}{}
			list = append(list, val)
		}
	}

	return list
}

func subtract(a []int, b []int) []int {
	store := make(map[int]interface{})
	for _, val := range b {
		store[val] = struct{}{}
	}

	var list []int
	for _, val := range a {
		if _, ok := store[val]; !ok {
			list = append(list, val)
		}
	}

	return list
}
