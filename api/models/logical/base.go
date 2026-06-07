package logical

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/objx"
)

type Relationship struct {
	SourcePath string   `json:"-"`
	Type       string   `json:"-"`
	Name       string   `json:"-"`
	ToOneId    string   `json:"-"`
	ToManyIds  []string `json:"-"`
	IsUUID     bool     `json:"-"`
}

// Fill the relationships IDs dynamically in ToOneId or ToManyIds
var ContextRelationships = []Relationship{
	{
		SourcePath: "Context.CreatedBy",
		Type:       "user",
		Name:       "createdBy",
	},
	{
		SourcePath: "Context.UpdatedBy",
		Type:       "user",
		Name:       "updatedBy",
	},
}

// ProjectAttributes defines model for ProjectAttributes.
type Base struct {
	Id            string
	Relationships []Relationship
	//Extend the derived model by adding its attributes
}

func (b *Base) SetToOneReferenceID(name, ID string) error {
	b.Relationships = append(b.Relationships, Relationship{Name: name, ToOneId: ID})
	return nil
}

func getJsonAPIReferences(relationships []Relationship) []jsonapi.Reference {
	var references []jsonapi.Reference
	for _, relation := range relationships {
		references = append(references, jsonapi.Reference{Type: relation.Type, Name: relation.Name})
	}

	return references
}

func getJsonAPIReferenceIDs(relationships []Relationship) []jsonapi.ReferenceID {
	var references []jsonapi.ReferenceID
	for _, relation := range relationships {
		if relation.ToOneId != "" {
			references = append(references, jsonapi.ReferenceID{
				Type: relation.Type,
				Name: relation.Name,
				ID:   relation.ToOneId,
			})
		}

		if relation.ToManyIds != nil {
			for _, Id := range relation.ToManyIds {
				references = append(references, jsonapi.ReferenceID{
					Type: relation.Type,
					Name: relation.Name,
					ID:   Id,
				})
			}
		}
	}

	return references
}

type UUIDInterface interface {
	String() string
}

func serializeRelationships(input interface{}, modelRelationships []Relationship) []Relationship {
	inputMap := map[string]interface{}{}

	//Convert input to reflection(Go reflect library) interface
	value := reflect.ValueOf(input)

	//Indirect the input value : This will allow both pointer or value as input
	indirectValue := reflect.Indirect(value)

	//Take action only if the input is of kind struct or map[string] interface{}
	kindOfInput := indirectValue.Kind()

	switch kindOfInput {
	case reflect.Struct:
		s := structs.New(input)
		inputMap = s.Map()
	case reflect.Map:
		inputMap = input.(map[string]interface{})
	}

	var result []Relationship
	if len(inputMap) > 0 {
		m := objx.New(inputMap)
		relationships := append(ContextRelationships, modelRelationships...)
		for index, relation := range relationships {
			if relation.SourcePath != "" {
				valueInPath := m.Get(relation.SourcePath)

				if !valueInPath.IsNil() {
					value := reflect.ValueOf(valueInPath.Data())

					kindOfRelation := value.Kind()

					if relation.IsUUID {
						relationships[index].ToOneId = valueInPath.Data().(UUIDInterface).String()
					} else {
						switch kindOfRelation {
						case reflect.Slice:
							if value.Len() > 0 {
								if value.Elem().String() == "string" {
									relationships[index].ToManyIds = valueInPath.MustStrSlice()
								}
							}
						case reflect.Array:
							if value.Len() > 0 {
								relationships[index].ToManyIds = valueInPath.MustStrSlice()
							}
						case reflect.String:
							if valueInPath.String() != "" {
								relationships[index].ToOneId = valueInPath.String()
							}
						case reflect.Ptr:
							indirectValue := reflect.Indirect(value)
							if indirectValue.Kind() == reflect.String {
								fieldValue := valueInPath.Data().(*string)
								if fieldValue != nil {
									relationships[index].ToOneId = *fieldValue
								}
							}
						}

					}

					if relationships[index].ToOneId != "" || relationships[index].ToManyIds != nil {
						result = append(result, relationships[index])
					}

				}
			}
		}
	}

	return result
}
