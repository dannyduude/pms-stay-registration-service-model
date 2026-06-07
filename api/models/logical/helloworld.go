package logical

import (
	"go-printos-backend-quickstart/api/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/manyminds/api2go/jsonapi"
)

const helloWorldType = "helloworld"

/******************************Version Response Model*************************/

// VersionAttributes defines model for VersionAttributes.

type HelloWorld struct {
	Base `json:"~"`
	models.HelloWorldAttributes
}

var helloWorldRelationships = []Relationship{}

func (v HelloWorld) GetID() string {
	return v.Id
}

func (v *HelloWorld) SetID(id string) error {
	//No action as ID is not set as part of request
	return nil
}

func (v HelloWorld) GetName() string {
	return helloWorldType
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (v HelloWorld) GetReferences() []jsonapi.Reference {
	references := getJsonAPIReferences(v.Relationships)
	return references
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (v HelloWorld) GetReferencedIDs() []jsonapi.ReferenceID {
	result := getJsonAPIReferenceIDs(v.Relationships)
	return result
}

func (v *HelloWorld) SerializeRelationships(input interface{}) {
	v.Relationships = serializeRelationships(input, helloWorldRelationships)
}

/*************************************************************************************/

/******************************HelloWorld POST Model**********************************/

type POSTHelloWorldAttributes models.POSTHelloWorldAttributes

func (v *POSTHelloWorldAttributes) SetID(id string) error {
	//No action as ID is not set as part of request
	return nil
}

func (v POSTHelloWorldAttributes) GetName() string {
	return helloWorldType
}

func (p POSTHelloWorldAttributes) ValidateSchema() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Message, validation.Required),
	)
}

func (p POSTHelloWorldAttributes) ValidateData() error {
	return nil
}

/*************************************************************************************/

// /******************************Project PATCH Model*************************/

// type PATCHHelloWorldAttributes models.PATCHVersionAttributes

// func (p *PATCHVersionAttributes) SetID(id string) error {
// 	//No action as ID is not set as part of request
// 	return nil
// }

// func (p PATCHVersionAttributes) GetName() string {
// 	return versionType
// }

// func (p PATCHVersionAttributes) ValidateSchema() error {
// 	return nil // no fields are required
// }

// func (p PATCHVersionAttributes) ValidateData() error {
// 	if p.Processes != nil {
// 		return validation.ValidateStruct(&p,
// 			validation.Field(&p.Processes, validation.By(func(value interface{}) error {
// 				processes := *value.(*[]models.ProcessesArrayItem)
// 				var err error
// 				for i := range processes {
// 					if err = validates.ProcessesArrayItem(processes[i]); err != nil {
// 						break
// 					}
// 				}

// 				return err
// 			})),
// 		)
// 	}

// 	return nil
// }

// /*************************************************************************************/
// /******************************Post Version Submit Model*************************/
// type POSTVersionSubmitAttributes models.POSTVersionSubmitAttributes

// func (v POSTVersionSubmitAttributes) GetID() string {
// 	return ""
// }
// func (v POSTVersionSubmitAttributes) SetID(id string) error {
// 	//No action as ID is not set as part of request
// 	return nil
// }

// func (p POSTVersionSubmitAttributes) GetName() string {
// 	return jobSubmissionType
// }

// /*************************************************************************************/
