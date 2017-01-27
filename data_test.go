package bahamut

import "fmt"
import "github.com/aporeto-inc/elemental"

// ListIdentity represents the Identity of the object
var ListIdentity = elemental.Identity{
	Name:     "list",
	Category: "lists",
}

// ListsList represents a list of Lists
type ListsList []*List

// List represents the model of a list
type List struct {
	// The identifier
	ID string `json:"ID" cql:"id,omitempty" bson:"id"`

	// A creation only only attribute
	CreationOnly string `json:"creationOnly" cql:"creationonly,omitempty" bson:"creationonly"`

	// The description
	Description string `json:"description" cql:"description,omitempty" bson:"description"`

	// The name
	Name string `json:"name" cql:"name,omitempty" bson:"name"`

	// The identifier of the parent of the object
	ParentID string `json:"parentID" cql:"parentid,omitempty" bson:"parentid"`

	// The type of the parent of the object
	ParentType string `json:"parentType" cql:"parenttype,omitempty" bson:"parenttype"`

	// A read only attribute
	ReadOnly string `json:"readOnly" cql:"readonly,omitempty" bson:"readonly"`
}

// NewList returns a new *List
func NewList() *List {

	return &List{}
}

// Identity returns the Identity of the object.
func (o *List) Identity() elemental.Identity {

	return ListIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *List) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *List) SetIdentifier(ID string) {

	o.ID = ID
}

func (o *List) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetCreationOnly returns the creationOnly of the receiver
func (o *List) GetCreationOnly() string {
	return o.CreationOnly
}

// GetName returns the name of the receiver
func (o *List) GetName() string {
	return o.Name
}

// SetName set the given name of the receiver
func (o *List) SetName(name string) {
	o.Name = name
}

// GetReadOnly returns the readOnly of the receiver
func (o *List) GetReadOnly() string {
	return o.ReadOnly
}

// Validate valides the current information stored into the structure.
func (o *List) Validate() error {

	errors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("creationOnly", o.CreationOnly); err != nil {
		errors = append(errors, err)
	}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		errors = append(errors, err)
	}

	if err := elemental.ValidateRequiredString("readOnly", o.ReadOnly); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// SpecificationForAttribute returns the AttributeSpecification for the given attribute name key.
func (List) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	return ListAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (List) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return ListAttributesMap
}

// ListAttributesMap represents the map of attribute for List.
var ListAttributesMap = map[string]elemental.AttributeSpecification{
	"ID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Identifier:     true,
		Name:           "ID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"CreationOnly": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		CreationOnly:   true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Getter:         true,
		Name:           "creationOnly",
		Orderable:      true,
		Required:       true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"Description": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "description",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
	},
	"Name": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Getter:         true,
		Name:           "name",
		Orderable:      true,
		Required:       true,
		Setter:         true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ParentID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		ForeignKey:     true,
		Format:         "free",
		Name:           "parentID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ParentType": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "parentType",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ReadOnly": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Getter:         true,
		Name:           "readOnly",
		Orderable:      true,
		ReadOnly:       true,
		Required:       true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
}

// TaskStatusValue represents the possible values for attribute "status".
type TaskStatusValue string

const (
	// TaskStatusDone represents the value DONE.
	TaskStatusDone TaskStatusValue = "DONE"

	// TaskStatusProgress represents the value PROGRESS.
	TaskStatusProgress TaskStatusValue = "PROGRESS"

	// TaskStatusTodo represents the value TODO.
	TaskStatusTodo TaskStatusValue = "TODO"
)

// TaskIdentity represents the Identity of the object
var TaskIdentity = elemental.Identity{
	Name:     "task",
	Category: "tasks",
}

// TasksList represents a list of Tasks
type TasksList []*Task

// Task represents the model of a task
type Task struct {
	// The identifier
	ID string `json:"ID" cql:"id,omitempty" bson:"id"`

	// The description
	Description string `json:"description" cql:"description,omitempty" bson:"description"`

	// The name
	Name string `json:"name" cql:"name,omitempty" bson:"name"`

	// The identifier of the parent of the object
	ParentID string `json:"parentID" cql:"parentid,omitempty" bson:"parentid"`

	// The type of the parent of the object
	ParentType string `json:"parentType" cql:"parenttype,omitempty" bson:"parenttype"`

	// The status of the task
	Status TaskStatusValue `json:"status" cql:"status,omitempty" bson:"status"`
}

// NewTask returns a new *Task
func NewTask() *Task {

	return &Task{
		Status: "TODO",
	}
}

// Identity returns the Identity of the object.
func (o *Task) Identity() elemental.Identity {

	return TaskIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Task) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Task) SetIdentifier(ID string) {

	o.ID = ID
}

func (o *Task) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// Validate valides the current information stored into the structure.
func (o *Task) Validate() error {

	errors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		errors = append(errors, err)
	}

	if err := elemental.ValidateStringInList("status", string(o.Status), []string{"DONE", "PROGRESS", "TODO"}, false); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// SpecificationForAttribute returns the AttributeSpecification for the given attribute name key.
func (Task) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	return TaskAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (Task) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return TaskAttributesMap
}

// TaskAttributesMap represents the map of attribute for Task.
var TaskAttributesMap = map[string]elemental.AttributeSpecification{
	"ID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Identifier:     true,
		Name:           "ID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"Description": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "description",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
	},
	"Name": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "name",
		Orderable:      true,
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"ParentID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		ForeignKey:     true,
		Format:         "free",
		Name:           "parentID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ParentType": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "parentType",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"Status": elemental.AttributeSpecification{
		AllowedChoices: []string{"DONE", "PROGRESS", "TODO"},
		Exposed:        true,
		Filterable:     true,
		Name:           "status",
		Orderable:      true,
		Stored:         true,
		Type:           "enum",
	},
}
var RootIdentity = elemental.Identity{
	Name:     "root",
	Category: "root",
}

// Root represents the model of a root
type Root struct {
	// The identifier
	ID string `json:"ID" cql:"id,omitempty" bson:"id"`

	// The identifier of the parent of the object
	ParentID string `json:"parentID" cql:"parentid,omitempty" bson:"parentid"`

	// The type of the parent of the object
	ParentType string `json:"parentType" cql:"parenttype,omitempty" bson:"parenttype"`

	Token        string `json:"APIKey,omitempty"`
	Organization string `json:"enterprise,omitempty"`
}

// NewRoot returns a new *Root
func NewRoot() *Root {

	return &Root{}
}

// Identity returns the Identity of the object.
func (o *Root) Identity() elemental.Identity {

	return RootIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Root) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Root) SetIdentifier(ID string) {

	o.ID = ID
}

func (o *Root) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// Validate valides the current information stored into the structure.
func (o *Root) Validate() error {

	errors := elemental.Errors{}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// APIKey returns a the API Key
func (o *Root) APIKey() string {

	return o.Token
}

// SetAPIKey sets a the API Key
func (o *Root) SetAPIKey(key string) {

	o.Token = key
}

// SpecificationForAttribute returns the AttributeSpecification for the given attribute name key.
func (Root) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	return RootAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (Root) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return RootAttributesMap
}

// RootAttributesMap represents the map of attribute for Root.
var RootAttributesMap = map[string]elemental.AttributeSpecification{
	"ID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Identifier:     true,
		Name:           "ID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ParentID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		ForeignKey:     true,
		Format:         "free",
		Name:           "parentID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ParentType": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "parentType",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
}

// UserIdentity represents the Identity of the object
var UserIdentity = elemental.Identity{
	Name:     "user",
	Category: "users",
}

// UsersList represents a list of Users
type UsersList []*User

// User represents the model of a user
type User struct {
	// The identifier
	ID string `json:"ID" cql:"id,omitempty" bson:"id"`

	// The first name
	FirstName string `json:"firstName" cql:"firstname,omitempty" bson:"firstname"`

	// The last name
	LastName string `json:"lastName" cql:"lastname,omitempty" bson:"lastname"`

	// The identifier of the parent of the object
	ParentID string `json:"parentID" cql:"parentid,omitempty" bson:"parentid"`

	// The type of the parent of the object
	ParentType string `json:"parentType" cql:"parenttype,omitempty" bson:"parenttype"`

	// the login
	UserName string `json:"userName" cql:"username,omitempty" bson:"username"`
}

// NewUser returns a new *User
func NewUser() *User {

	return &User{}
}

// Identity returns the Identity of the object.
func (o *User) Identity() elemental.Identity {

	return UserIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *User) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *User) SetIdentifier(ID string) {

	o.ID = ID
}

func (o *User) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// Validate valides the current information stored into the structure.
func (o *User) Validate() error {

	errors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("firstName", o.FirstName); err != nil {
		errors = append(errors, err)
	}

	if err := elemental.ValidateRequiredString("lastName", o.LastName); err != nil {
		errors = append(errors, err)
	}

	if err := elemental.ValidateRequiredString("userName", o.UserName); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// SpecificationForAttribute returns the AttributeSpecification for the given attribute name key.
func (User) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	return UserAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (User) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return UserAttributesMap
}

// UserAttributesMap represents the map of attribute for User.
var UserAttributesMap = map[string]elemental.AttributeSpecification{
	"ID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Identifier:     true,
		Name:           "ID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"FirstName": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "firstName",
		Orderable:      true,
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"LastName": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "lastName",
		Orderable:      true,
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"ParentID": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		ForeignKey:     true,
		Format:         "free",
		Name:           "parentID",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"ParentType": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Autogenerated:  true,
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "parentType",
		Orderable:      true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
	"UserName": elemental.AttributeSpecification{
		AllowedChoices: []string{},
		Exposed:        true,
		Filterable:     true,
		Format:         "free",
		Name:           "userName",
		Orderable:      true,
		Required:       true,
		Stored:         true,
		Type:           "string",
		Unique:         true,
	},
}

var relationshipsRegistry elemental.RelationshipsRegistry

// Relationships returns the model relationships.
func Relationships() elemental.RelationshipsRegistry {

	return relationshipsRegistry
}

func init() {
	relationshipsRegistry = elemental.RelationshipsRegistry{}

	//
	// Main Relationship for root
	//
	RootMainRelationship := &elemental.Relationship{
		AllowsRetrieve: true,
	}

	// Children relationship for lists in root
	RootMainRelationship.AddChild(
		elemental.IdentityFromName("list"),
		&elemental.Relationship{
			AllowsCreate:       true,
			AllowsRetrieveMany: true,
			AllowsInfo:         true,
		},
	)
	// Children relationship for users in root
	RootMainRelationship.AddChild(
		elemental.IdentityFromName("user"),
		&elemental.Relationship{
			AllowsCreate:       true,
			AllowsRetrieveMany: true,
			AllowsInfo:         true,
		},
	)

	relationshipsRegistry[elemental.IdentityFromName("root")] = RootMainRelationship

	//
	// Main Relationship for task
	//
	TaskMainRelationship := &elemental.Relationship{
		AllowsRetrieve: true,
		AllowsUpdate:   true,
		AllowsDelete:   true,
	}

	relationshipsRegistry[elemental.IdentityFromName("task")] = TaskMainRelationship

	//
	// Main Relationship for list
	//
	ListMainRelationship := &elemental.Relationship{
		AllowsRetrieve: true,
		AllowsUpdate:   true,
		AllowsDelete:   true,
	}

	// Children relationship for tasks in list
	ListMainRelationship.AddChild(
		elemental.IdentityFromName("task"),
		&elemental.Relationship{
			AllowsCreate:       true,
			AllowsRetrieveMany: true,
			AllowsInfo:         true,
		},
	)
	// Children relationship for users in list
	ListMainRelationship.AddChild(
		elemental.IdentityFromName("user"),
		&elemental.Relationship{
			AllowsPatch:        true,
			AllowsRetrieveMany: true,
			AllowsInfo:         true,
		},
	)

	relationshipsRegistry[elemental.IdentityFromName("list")] = ListMainRelationship

	//
	// Main Relationship for user
	//
	UserMainRelationship := &elemental.Relationship{
		AllowsRetrieve: true,
		AllowsUpdate:   true,
		AllowsDelete:   true,
	}

	relationshipsRegistry[elemental.IdentityFromName("user")] = UserMainRelationship

}

func init() {

	elemental.RegisterIdentity(RootIdentity)
	elemental.RegisterIdentity(TaskIdentity)
	elemental.RegisterIdentity(ListIdentity)
	elemental.RegisterIdentity(UserIdentity)
}

// IdentifiableForIdentity returns a new instance of the Identifiable for the given identity name.
func IdentifiableForIdentity(identity string) elemental.Identifiable {

	switch identity {
	case RootIdentity.Name:
		return NewRoot()
	case TaskIdentity.Name:
		return NewTask()
	case ListIdentity.Name:
		return NewList()
	case UserIdentity.Name:
		return NewUser()
	default:
		return nil
	}
}

var UnmarshalableListIdentity = elemental.Identity{Name: "list", Category: "lists"}

type UnmarshalableList struct {
	List
}

func NewUnmarshalableList() *UnmarshalableList {
	return &UnmarshalableList{List: List{}}
}

func (o *UnmarshalableList) Identity() elemental.Identity { return UnmarshalableListIdentity }

func (o *UnmarshalableList) UnmarshalJSON([]byte) error {
	return fmt.Errorf("error unmarshalling")
}

func (o *UnmarshalableList) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("error marshalling")
}

func (o *UnmarshalableList) Validate() elemental.Errors { return nil }