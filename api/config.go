package api

type KanataConfig struct {
	NodeId uint32

	StoreProvider   string
	UUIDProvider    string
	MarshalProvider string

	StoreConfig *StoreInitConfig
}

type KanataConfigParser interface {
	ParseConfigFile(path string) (*KanataConfig, error)
}

var (
	kanataConfigParserProviders = make(map[string]KanataConfigParser)
)

func GetKanataConfigParser(name string) KanataConfigParser {
	p, ok := kanataConfigParserProviders[name]
	if !ok {
		return nil
	} else {
		return p
	}
}

func RegisterKanataConfigParse(name string, parser KanataConfigParser) {
	kanataConfigParserProviders[name] = parser
}
