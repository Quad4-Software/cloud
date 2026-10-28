package identity

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/opencloud-eu/libre-graph-api-go"
	"github.com/Quad4-Software/cloud/services/graph/pkg/config"
	"github.com/Quad4-Software/cloud/services/graph/pkg/errorcode"
	"github.com/Quad4-Software/cloud/services/graph/pkg/identity/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var eduConfig = config.LDAP{
	UserBaseDN:               "ou=people,dc=test",
	UserObjectClass:          "inetOrgPerson",
	UserSearchScope:          "sub",
	UserFilter:               "",
	UserDisplayNameAttribute: "displayname",
	UserIDAttribute:          "entryUUID",
	UserEmailAttribute:       "mail",
	UserNameAttribute:        "uid",
	UserEnabledAttribute:     "userEnabledAttribute",
	DisableUserMechanism:     "attribute",
	UserTypeAttribute:        "userTypeAttribute",

	GroupBaseDN:          "ou=groups,dc=test",
	GroupObjectClass:     "groupOfNames",
	GroupSearchScope:     "sub",
	GroupFilter:          "",
	GroupNameAttribute:   "cn",
	GroupMemberAttribute: "member",
	GroupIDAttribute:     "entryUUID",

	WriteEnabled:              true,
	EducationResourcesEnabled: true,
}

var schoolEntry = ldap.NewEntry("ou=Test School",
	map[string][]string{
		"ou":                             {"Test School"},
		"quad4CloudEducationSchoolNumber": {"0123"},
		"quad4CloudUUID":                  {"abcd-defg"},
		"quad4CloudEducationExternalId":   {"abcd-defg"},
	})

var schoolEntry1 = ldap.NewEntry("ou=Test School1",
	map[string][]string{
		"ou":                             {"Test School1"},
		"quad4CloudEducationSchoolNumber": {"0042"},
		"quad4CloudUUID":                  {"hijk-defg"},
		"quad4CloudEducationExternalId":   {"hijk-defg"},
	})
var schoolEntryWithTermination = ldap.NewEntry("ou=Test School",
	map[string][]string{
		"ou":                             {"Test School"},
		"quad4CloudEducationSchoolNumber": {"0123"},
		"quad4CloudUUID":                  {"abcd-defg"},
		"quad4CloudEducationExternalId":   {"abcd-defg"},
		"quad4CloudEducationSchoolTerminationTimestamp": {"20420131120000Z"},
	})

var (
	filterSchoolSearchByIdExisting        = "(&(objectClass=quad4CloudEducationSchool)(|(quad4CloudUUID=abcd-defg)(quad4CloudEducationSchoolNumber=abcd-defg)))"
	filterSchoolSearchByIdNonexistant     = "(&(objectClass=quad4CloudEducationSchool)(|(quad4CloudUUID=xxxx-xxxx)(quad4CloudEducationSchoolNumber=xxxx-xxxx)))"
	filterSchoolSearchByNumberExisting    = "(&(objectClass=quad4CloudEducationSchool)(|(quad4CloudUUID=0123)(quad4CloudEducationSchoolNumber=0123)))"
	filterSchoolSearchByNumberNonexistant = "(&(objectClass=quad4CloudEducationSchool)(|(quad4CloudUUID=3210)(quad4CloudEducationSchoolNumber=3210)))"

	schoolLDAPAttributeTypes = []string{"ou", "quad4CloudUUID", "quad4CloudEducationExternalId", "quad4CloudEducationSchoolNumber", "quad4CloudEducationSchoolTerminationTimestamp"}
)

func TestCreateEducationSchool(t *testing.T) {
	tests := []struct {
		name          string
		schoolNumber  string
		schoolName    string
		expectedError error
	}{
		{
			name:          "Create a Education School succeeds",
			schoolNumber:  "0123",
			schoolName:    "Test School",
			expectedError: nil,
		}, {
			name:          "Create a Education School with a duplicated Schoolnumber fails with an error",
			schoolNumber:  "0666",
			schoolName:    "Test School",
			expectedError: errorcode.New(errorcode.NameAlreadyExists, "A school with that number is already present"),
		}, {
			name:          "Create a Education School with a duplicated Name fails with an error",
			schoolNumber:  "0123",
			schoolName:    "Existing Test School",
			expectedError: errorcode.New(errorcode.NameAlreadyExists, "A school with that name is already present"),
		}, {
			name:          "Create a Education School fails, when there is a backend error",
			schoolNumber:  "1111",
			schoolName:    "Test School",
			expectedError: errorcode.New(errorcode.GeneralException, "error looking up school by number"),
		},
	}
	for _, tt := range tests {
		lm := &mocks.Client{}
		ldapSchoolGoodAddRequestMatcher := func(ar *ldap.AddRequest) bool {
			if ar.DN != "ou=Test School," {
				return false
			}
			for _, attr := range ar.Attributes {
				if attr.Type == "quad4CloudEducationSchoolTerminationTimestamp" {
					return false
				}
			}
			return true
		}
		lm.On("Add", mock.MatchedBy(ldapSchoolGoodAddRequestMatcher)).Return(nil)

		ldapExistingSchoolAddRequestMatcher := func(ar *ldap.AddRequest) bool {
			if ar.DN == "ou=Existing Test School," {
				return true
			}
			return false
		}
		lm.On("Add", mock.MatchedBy(ldapExistingSchoolAddRequestMatcher)).Return(ldap.NewError(ldap.LDAPResultEntryAlreadyExists, errors.New("")))

		schoolNumberSearchRequest := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     "(&(objectClass=quad4CloudEducationSchool)(quad4CloudEducationSchoolNumber=0123))",
			Attributes: schoolLDAPAttributeTypes,
			Controls:   []ldap.Control(nil),
		}
		lm.On("Search", schoolNumberSearchRequest).
			Return(
				&ldap.SearchResult{
					Entries: []*ldap.Entry{},
				},
				nil)
		existingSchoolNumberSearchRequest := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     "(&(objectClass=quad4CloudEducationSchool)(quad4CloudEducationSchoolNumber=0666))",
			Attributes: schoolLDAPAttributeTypes,
			Controls:   []ldap.Control(nil),
		}
		lm.On("Search", existingSchoolNumberSearchRequest).
			Return(
				&ldap.SearchResult{
					Entries: []*ldap.Entry{schoolEntry},
				},
				nil)
		schoolNumberSearchRequestError := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     "(&(objectClass=quad4CloudEducationSchool)(quad4CloudEducationSchoolNumber=1111))",
			Attributes: schoolLDAPAttributeTypes,
			Controls:   []ldap.Control(nil),
		}
		lm.On("Search", schoolNumberSearchRequestError).
			Return(
				&ldap.SearchResult{
					Entries: []*ldap.Entry{},
				},
				ldap.NewError(ldap.LDAPResultOther, errors.New("some error")))
		schoolLookupAfterCreate := &ldap.SearchRequest{
			BaseDN:     "ou=Test School,",
			Scope:      0,
			SizeLimit:  1,
			Filter:     "(objectClass=quad4CloudEducationSchool)",
			Attributes: schoolLDAPAttributeTypes,
			Controls:   []ldap.Control(nil),
		}
		lm.On("Search", schoolLookupAfterCreate).
			Return(
				&ldap.SearchResult{
					Entries: []*ldap.Entry{schoolEntry},
				},
				nil)
		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)
		assert.NotEqual(t, "", b.educationConfig.schoolObjectClass)

		school := libregraph.NewEducationSchool()
		school.SetDisplayName(tt.schoolName)
		school.SetSchoolNumber(tt.schoolNumber)
		school.SetId("abcd-defg")
		resSchool, err := b.CreateEducationSchool(context.Background(), *school)
		if tt.expectedError == nil {
			assert.Nil(t, err)
			lm.AssertNumberOfCalls(t, "Add", 1)
			assert.NotNil(t, resSchool)
			assert.Equal(t, resSchool.GetDisplayName(), school.GetDisplayName())
			assert.Equal(t, resSchool.GetId(), school.GetId())
			assert.Equal(t, resSchool.GetSchoolNumber(), school.GetSchoolNumber())
			assert.False(t, resSchool.HasTerminationDate())
		} else {
			assert.Equal(t, err, tt.expectedError)
			assert.Nil(t, resSchool)
		}
	}
}

func TestUpdateEducationSchoolTerminationDate(t *testing.T) {
	lm := &mocks.Client{}

	ldapSchoolTerminationRequestMatcher := func(mr *ldap.ModifyRequest) bool {
		if mr.DN != "ou=Test School" {
			return false
		}
		for _, mod := range mr.Changes {
			if mod.Operation == ldap.ReplaceAttribute &&
				mod.Modification.Type == "quad4CloudEducationSchoolTerminationTimestamp" &&
				mod.Modification.Vals[0] == "20420131120000Z" {
				return true
			}
		}
		return false
	}
	lm.On("Modify", mock.MatchedBy(ldapSchoolTerminationRequestMatcher)).Return(nil)
	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{schoolEntry},
			},
			nil).
		Once()
	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{schoolEntryWithTermination},
			},
			nil).
		Once()

	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	assert.NotEqual(t, "", b.educationConfig.schoolObjectClass)
	school := libregraph.NewEducationSchool()
	terminationTime := time.Date(2042, time.January, 31, 12, 0, 0, 0, time.UTC)
	school.SetTerminationDate(terminationTime)
	resSchool, err := b.UpdateEducationSchool(context.Background(), "abcd-defg", *school)
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.Nil(t, err)
	assert.NotNil(t, resSchool)
	assert.Equal(t, "Test School", resSchool.GetDisplayName())
	assert.Equal(t, "abcd-defg", resSchool.GetId())
	assert.Equal(t, "0123", resSchool.GetSchoolNumber())
	assert.True(t, resSchool.HasTerminationDate())
	assert.True(t, terminationTime.Equal(resSchool.GetTerminationDate()))
}

func TestUpdateEducationSchoolOperation(t *testing.T) {
	testSchoolName := "A name"
	testSchoolNumber := "1234"
	tests := []struct {
		name              string
		displayName       string
		schoolNumber      string
		expectedOperation schoolUpdateOperation
	}{
		{
			name:              "Test using school with both number and name, unchanged",
			displayName:       testSchoolName,
			schoolNumber:      testSchoolNumber,
			expectedOperation: schoolUnchanged,
		},
		{
			name:              "Test using school with both number and name, unchanged",
			displayName:       "A new name",
			schoolNumber:      "9876",
			expectedOperation: tooManyValues,
		},
		{
			name:              "Test with unchanged number",
			schoolNumber:      testSchoolNumber,
			expectedOperation: schoolUnchanged,
		},
		{
			name:              "Test with unchanged name",
			displayName:       testSchoolName,
			expectedOperation: schoolUnchanged,
		},
		{
			name:              "Test new name",
			displayName:       "Something new",
			expectedOperation: schoolRenamed,
		},
		{
			name:              "Test new number",
			schoolNumber:      "9876",
			expectedOperation: schoolPropertiesUpdated,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		displayName := "A name"
		schoolNumber := "1234"

		currentSchool := libregraph.EducationSchool{
			DisplayName:  &displayName,
			SchoolNumber: &schoolNumber,
		}

		schoolUpdate := libregraph.EducationSchool{
			DisplayName:  &tt.displayName,
			SchoolNumber: &tt.schoolNumber,
		}

		operation := b.updateEducationSchoolOperation(schoolUpdate, currentSchool)
		assert.Equal(t, tt.expectedOperation, operation)
	}
}

func TestDeleteEducationSchool(t *testing.T) {
	tests := []struct {
		name                 string
		numberOrId           string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test delete school using schoolId",
			numberOrId:           "abcd-defg",
			filter:               filterSchoolSearchByIdExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test delete school using unknown schoolId",
			numberOrId:           "xxxx-xxxx",
			filter:               filterSchoolSearchByIdNonexistant,
			expectedItemNotFound: true,
		},
		{
			name:                 "Test delete school using schoolNumber",
			numberOrId:           "0123",
			filter:               filterSchoolSearchByNumberExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test delete school using unknown schoolNumber",
			numberOrId:           "3210",
			filter:               filterSchoolSearchByNumberNonexistant,
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		sr := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: schoolLDAPAttributeTypes,
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
		}
		dr := &ldap.DelRequest{
			DN: "ou=Test School",
		}
		lm.On("Del", dr).Return(nil)

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		err = b.DeleteEducationSchool(context.Background(), tt.numberOrId)
		lm.AssertNumberOfCalls(t, "Search", 1)

		if tt.expectedItemNotFound {
			lm.AssertNumberOfCalls(t, "Del", 0)
			assert.NotNil(t, err)
			assert.Equal(t, "itemNotFound: not found", err.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetEducationSchool(t *testing.T) {
	tests := []struct {
		name                 string
		numberOrId           string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test search school using schoolId",
			numberOrId:           "abcd-defg",
			filter:               filterSchoolSearchByIdExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown schoolId",
			numberOrId:           "xxxx-xxxx",
			filter:               filterSchoolSearchByIdNonexistant,
			expectedItemNotFound: true,
		},
		{
			name:                 "Test search school using schoolNumber",
			numberOrId:           "0123",
			filter:               filterSchoolSearchByNumberExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown schoolNumber",
			numberOrId:           "3210",
			filter:               filterSchoolSearchByNumberNonexistant,
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		sr := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: schoolLDAPAttributeTypes,
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
		}

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		school, err := b.GetEducationSchool(context.Background(), tt.numberOrId)
		lm.AssertNumberOfCalls(t, "Search", 1)

		if tt.expectedItemNotFound {
			assert.NotNil(t, err)
			assert.Equal(t, "itemNotFound: not found", err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, "Test School", school.GetDisplayName())
			assert.Equal(t, "abcd-defg", school.GetId())
			assert.Equal(t, "0123", school.GetSchoolNumber())
		}
	}
}

func TestGetEducationSchools(t *testing.T) {
	lm := &mocks.Client{}
	sr1 := &ldap.SearchRequest{
		BaseDN:     "",
		Scope:      2,
		SizeLimit:  0,
		Filter:     "(objectClass=quad4CloudEducationSchool)",
		Attributes: schoolLDAPAttributeTypes,
		Controls:   []ldap.Control(nil),
	}
	lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry, schoolEntry1}}, nil)
	//	lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	_, err = b.GetEducationSchools(context.Background())
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
}

var schoolByIDSearch1 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "",
	Scope:      2,
	SizeLimit:  1,
	Filter:     filterSchoolSearchByIdExisting,
	Attributes: schoolLDAPAttributeTypes,
	Controls:   []ldap.Control(nil),
}

var schoolByNumberSearch *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "",
	Scope:      2,
	SizeLimit:  1,
	Filter:     filterSchoolSearchByNumberExisting,
	Attributes: schoolLDAPAttributeTypes,
	Controls:   []ldap.Control(nil),
}

var userByIDSearch1 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=quad4CloudEducationUser)(|(uid=abcd-defg)(entryUUID=abcd-defg)))",
	Attributes: eduUserAttrs,
	Controls:   []ldap.Control(nil),
}

var userByIDSearch2 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=quad4CloudEducationUser)(|(uid=does-not-exist)(entryUUID=does-not-exist)))",
	Attributes: eduUserAttrs,
	Controls:   []ldap.Control(nil),
}

var userToSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "uid=user,ou=people,dc=test",
	Changes: []ldap.Change{
		{
			Operation: ldap.AddAttribute,
			Modification: ldap.PartialAttribute{
				Type: "quad4CloudMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

var userFromSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "uid=user,ou=people,dc=test",
	Changes: []ldap.Change{
		{
			Operation: ldap.DeleteAttribute,
			Modification: ldap.PartialAttribute{
				Type: "quad4CloudMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

var classToSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "quad4CloudEducationExternalId=Math0123",
	Changes: []ldap.Change{
		{
			Operation: ldap.AddAttribute,
			Modification: ldap.PartialAttribute{
				Type: "quad4CloudMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

var classFromSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "quad4CloudEducationExternalId=Math0123",
	Changes: []ldap.Change{
		{
			Operation: ldap.DeleteAttribute,
			Modification: ldap.PartialAttribute{
				Type: "quad4CloudMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

func TestAddUsersToEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", userByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntry}}, nil)
	lm.On("Search", userByIDSearch2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Modify", userToSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.AddUsersToEducationSchool(context.Background(), "abcd-defg", []string{"does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	err = b.AddUsersToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg", "does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 5)
	assert.NotNil(t, err)
	err = b.AddUsersToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 7)
	assert.Nil(t, err)
	// try to add by school number (instead or id)
	err = b.AddUsersToEducationSchool(context.Background(), "0123", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 9)
	assert.Nil(t, err)
}

func TestRemoveMemberFromEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", userByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntryWithSchool}}, nil)
	lm.On("Search", userByIDSearch2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Modify", userFromSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.RemoveUserFromEducationSchool(context.Background(), "abcd-defg", "does-not-exist")
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound: not found", err.Error())
	err = b.RemoveUserFromEducationSchool(context.Background(), "abcd-defg", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 4)
	lm.AssertNumberOfCalls(t, "Modify", 1)
	// try to remove by school number (instead or id)
	err = b.RemoveUserFromEducationSchool(context.Background(), "0123", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 6)
	lm.AssertNumberOfCalls(t, "Modify", 2)
	assert.Nil(t, err)
}

var usersBySchoolIDSearch *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  0,
	Filter:     "(&(objectClass=quad4CloudEducationUser)(quad4CloudMemberOfSchool=abcd-defg))",
	Attributes: eduUserAttrs,
	Controls:   []ldap.Control(nil),
}

func TestGetEducationSchoolUsers(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", usersBySchoolIDSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntryWithSchool}}, nil)
	b, _ := getMockedBackend(lm, eduConfig, &logger)
	users, err := b.GetEducationSchoolUsers(context.Background(), "abcd-defg")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	users, err = b.GetEducationSchoolUsers(context.Background(), "0123")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

var classesBySchoolIDSearch *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=groups,dc=test",
	Scope:      2,
	SizeLimit:  0,
	Filter:     "(&(objectClass=quad4CloudEducationClass)(quad4CloudMemberOfSchool=abcd-defg))",
	Attributes: []string{"cn", "entryUUID", "quad4CloudEducationClassType", "quad4CloudEducationExternalId", "quad4CloudMemberOfSchool", "quad4CloudEducationTeacherMember"},
	Controls:   []ldap.Control(nil),
}

func TestGetEducationSchoolClasses(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", classesBySchoolIDSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
	b, _ := getMockedBackend(lm, eduConfig, &logger)
	users, err := b.GetEducationSchoolClasses(context.Background(), "abcd-defg")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	users, err = b.GetEducationSchoolClasses(context.Background(), "0123")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

var classesByUUIDSearchNotFound *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=groups,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=quad4CloudEducationClass)(|(entryUUID=does-not-exist)(quad4CloudEducationExternalId=does-not-exist)))",
	Attributes: []string{"cn", "entryUUID", "quad4CloudEducationClassType", "quad4CloudEducationExternalId", "quad4CloudMemberOfSchool", "quad4CloudEducationTeacherMember"},
	Controls:   []ldap.Control(nil),
}

var classesByUUIDSearchFound *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=groups,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=quad4CloudEducationClass)(|(entryUUID=abcd-defg)(quad4CloudEducationExternalId=abcd-defg)))",
	Attributes: []string{"cn", "entryUUID", "quad4CloudEducationClassType", "quad4CloudEducationExternalId", "quad4CloudMemberOfSchool", "quad4CloudEducationTeacherMember"},
	Controls:   []ldap.Control(nil),
}

func TestAddClassesToEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", classesByUUIDSearchNotFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Search", classesByUUIDSearchFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Modify", classToSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.AddClassesToEducationSchool(context.Background(), "abcd-defg", []string{"does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	err = b.AddClassesToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg", "does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 5)
	assert.NotNil(t, err)
	err = b.AddClassesToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 7)
	assert.Nil(t, err)
	// try to add by school number (instead or id)
	err = b.AddClassesToEducationSchool(context.Background(), "0123", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 9)
	assert.Nil(t, err)
}

func TestRemoveClassFromEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", classesByUUIDSearchFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntryWithSchool}}, nil)
	lm.On("Search", classesByUUIDSearchNotFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Modify", classFromSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.RemoveClassFromEducationSchool(context.Background(), "abcd-defg", "does-not-exist")
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound: not found", err.Error())
	err = b.RemoveClassFromEducationSchool(context.Background(), "abcd-defg", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 4)
	lm.AssertNumberOfCalls(t, "Modify", 1)
	// try to remove by school number (instead or id)
	err = b.RemoveClassFromEducationSchool(context.Background(), "0123", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 6)
	lm.AssertNumberOfCalls(t, "Modify", 2)
	assert.Nil(t, err)
}
