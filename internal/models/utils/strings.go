package utils

func ToStringArray(list []*string) []string {
	res := []string{}
	for _, v := range list {
		res = append(res, *v)
	}
	return res
}

func FromStringArray(list []string) []*string {
	res := []*string{}
	for i := range list {
		res = append(res, &list[i])
	}
	return res
}
