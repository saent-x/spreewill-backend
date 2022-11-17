package models

func (v *Vendor) CollectionName() string {
	return "vendors"
}

func (v *Customer) CollectionName() string {
	return "customers"
}

func (v *Post) CollectionName() string {
	return "posts"
}

func (v *Comment) CollectionName() string {
	return "comments"
}

func (v *Story) CollectionName() string {
	return "story"
}
