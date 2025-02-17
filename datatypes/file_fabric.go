package datatypes

type FabricModIdentifier struct {
	SchemaVersion int      `json:"schemaVersion"`
	Id            string   `json:"id"`
	Version       string   `json:"version"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Authors       []string `json:"authors"`
	Contact       struct {
		Homepage string `json:"homepage"`
		Issues   string `json:"issues"`
		Sources  string `json:"sources"`
	} `json:"contact"`
	License     string `json:"license"`
	Icon        string `json:"icon"`
	Environment string `json:"environment"`
	Entrypoints struct {
		Client []string `json:"client"`
		Server []string `json:"server"`
	} `json:"entrypoints"`
	Mixins        []string `json:"mixins"`
	AccessWidener string   `json:"accessWidener"`
	Depends       struct {
		Minecraft    string `json:"minecraft"`
		Fabricloader string `json:"fabricloader"`
		Java         string `json:"java"`
	} `json:"depends"`
}
