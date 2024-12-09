package consts

const (
	MimeTypeJPG  = `image/jpg`
	MimeTypeJPEG = `image/jpeg`
	MimeTypePNG  = `image/png`
)

var MimeTypeWithExt = map[string]string{
	MimeTypeJPG:  "jpg",
	MimeTypeJPEG: "jpeg",
	MimeTypePNG:  "png",
}

var (
	AllowedMimeTypesKTP    = []string{MimeTypeJPG, MimeTypeJPEG, MimeTypePNG}
	AllowedMimeTypesSelfie = []string{MimeTypeJPG, MimeTypeJPEG, MimeTypePNG}
)
