package utils

func ToStringPointer(str string) *string {
	v := &str
	return v
}

func ToBooleanPointer(b bool) *bool {
	v := &b
	return v
}
