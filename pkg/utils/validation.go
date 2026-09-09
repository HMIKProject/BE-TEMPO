package utils

import "github.com/go-playground/validator/v10"

// FormatValidationError menerjemahkan error gin validation menjadi kalimat rapi
func FormatValidationError(err error) []string {
	var errors []string

	// Pastikan error berasal dari validation
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// Jika bukan error validasi struktur, kembalikan pesan asli
		return []string{err.Error()}
	}

	for _, e := range validationErrs {
		switch e.Tag() {
		case "required":
			errors = append(errors, e.Field()+" wajib diisi")
		case "email":
			errors = append(errors, e.Field()+" harus berupa format email yang benar")
		case "min":
			errors = append(errors, e.Field()+" terlalu pendek")
		case "max":
			errors = append(errors, e.Field()+" terlalu panjang")
		default:
			errors = append(errors, e.Field()+" tidak valid ("+e.Tag()+")")
		}
	}
	return errors
}
