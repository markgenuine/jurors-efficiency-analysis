package main

// Len is part of sort.Interface.
func (d dataForSort) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d dataForSort) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d dataForSort) Less(i, j int) bool {
	return d[i].Efficiency > d[j].Efficiency
}
