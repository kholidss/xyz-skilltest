package presentation

type File struct {
	Name     string `json:"name"`
	Mimetype string `json:"mimetype"`
	Size     int    `json:"size"`
	File     []byte `json:"file"`
}
