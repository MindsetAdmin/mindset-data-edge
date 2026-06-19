package functions

// FunctionType définit le type de fonction
type FunctionType string

const (
	TypeConnector FunctionType = "connector"
	TypeTransform FunctionType = "transform"
	TypeCalculate FunctionType = "calculate"
	TypeCondition FunctionType = "condition"
	TypeOutput    FunctionType = "output"
)

// FunctionInfo est la structure exposée via l'API
type FunctionInfo struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // string, pas FunctionType (pour le JSON)
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Inputs      []Param `json:"inputs"`
	Outputs     []Param `json:"outputs"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color"`
}

// Param représente un paramètre d'entrée/sortie
type Param struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}

// IsConnector vérifie si le type est un connecteur
func (f *FunctionInfo) IsConnector() bool {
	return f.Type == "connector"
}

// IsTransform vérifie si le type est une transformation
func (f *FunctionInfo) IsTransform() bool {
	return f.Type == "transform"
}

// IsCalculate vérifie si le type est un calcul
func (f *FunctionInfo) IsCalculate() bool {
	return f.Type == "calculate"
}

// IsCondition vérifie si le type est une condition
func (f *FunctionInfo) IsCondition() bool {
	return f.Type == "condition"
}

// IsOutput vérifie si le type est une sortie
func (f *FunctionInfo) IsOutput() bool {
	return f.Type == "output"
}
