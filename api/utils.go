package api

var (
	uuidProviders = make(map[string]UUIDGenerator)
)

type UUIDGenerator interface {
	Generate() (string, error)
}

func GetUUIDProvider(name string) UUIDGenerator {
	p, ok := uuidProviders[name]
	if !ok {
		return nil
	} else {
		return p
	}
}

func RegisterUUIDProvider(name string, generator UUIDGenerator) {
	uuidProviders[name] = generator
}
