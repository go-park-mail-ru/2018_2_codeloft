package models

//easyjson:json
type Cell struct {
	Val string `json:"color"`
	//Mu  sync.Mutex `json:"-"`
}

type FieldSize Position

//easyjson:json
type FieldInfo struct {
	Size  FieldSize                       `json:"size"`
	Field [FIELD_HEIGHT][FIELD_WIDTH]Cell `json:"field"`
}
