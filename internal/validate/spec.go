package validate

import (
	"strings"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

// SpecValidator validates a swagger spec
type SpecValidator struct {
	schema       *spec.Schema // swagger 2.0 schema
	spec         *spec.Document
	KnownFormats strfmt.Registry
}

// NewSpecValidator creates a new swagger spec validator instance
func NewSpecValidator(schema *spec.Schema, formats strfmt.Registry) *SpecValidator {
	return &SpecValidator{
		schema:       schema,
		KnownFormats: formats,
	}
}

// Validate validates the swagger spec
func (s *SpecValidator) Validate(data interface{}) *Result {
	var sd *spec.Document

	switch v := data.(type) {
	case *spec.Document:
		sd = v
	}
	if sd == nil {
		return sErr(errors.New(500, "spec validator can only validate spec.Document objects"))
	}

	res := new(Result)
	schv := NewSchemaValidator(s.schema, nil, "", s.KnownFormats)
	res.Merge(schv.Validate(sd.Spec())) // -
	res.Merge(s.validateItems())
	res.Merge(s.validateUniqueSecurityScopes())
	res.Merge(s.validateUniqueScopesSecurityDefinitions())
	res.Merge(s.validatePathParamPresence())
	res.Merge(s.validateReferenced())
	res.Merge(s.validateRequiredDefinitions())
	res.Merge(s.validateParameters())
	res.Merge(s.validateReferencesValid())
	return res
}

func (s *SpecValidator) validateItems() *Result {
	// validate parameter, items, schema and response objects for presence of item if type is array
	return nil
}

func (s *SpecValidator) validateUniqueSecurityScopes() *Result {
	// Each authorization/security reference should contain only unique scopes.
	// (Example: For an oauth2 authorization/security requirement, when listing the required scopes,
	// each scope should only be listed once.)
	return nil
}

func (s *SpecValidator) validateUniqueScopesSecurityDefinitions() *Result {
	// Each authorization/security scope in an authorization/security definition should be unique.
	return nil
}

func (s *SpecValidator) validatePathParamPresence() *Result {
	// Each defined operation path parameters must correspond to a named element in the API's path pattern.
	// (For example, you cannot have a path parameter named id for the following path /pets/{petId} but you must have a path parameter named petId.)
	return nil
}

func (s *SpecValidator) validateReferenced() *Result {
	// Each referenceable definition must have references.
	return nil
}

func (s *SpecValidator) validateRequiredDefinitions() *Result {
	// Each definition property listed in the required array must be defined in the properties of the model
	return nil
}

func (s *SpecValidator) validateParameters() *Result {
	// each parameter should have a unique `name` and `type` combination
	// each operation should have only 1 parameter of type body
	// each api path should be non-verbatim (account for path param names) unique per method
	res := new(Result)
	knownPaths := make(map[string]string)
	for path, pi := range s.spec.Operations() {
		segments, params := parsePath(path)
		knowns := make([]string, 0, len(segments))
		for _, s := range segments {
			knowns = append(knowns, s)
		}
		for _, i := range params {
			knowns[i] = "!"
		}
		knownPath := strings.Join(knowns, "/")
		if orig, ok := knownPaths[knownPath]; ok {
			res.AddErrors(errors.New(422, "path %s overlaps with %s", path, orig))
		} else {
			knownPaths[knownPath] = path
		}

		for _, op := range pi {
			ptypes := make(map[string]map[string]struct{})
			var firstBodyParam string

			for _, pr := range op.Parameters {
				if pr.In == "body" {
					if firstBodyParam != "" {
						res.AddErrors(errors.New(422, "operation %q has more than 1 body param (accepted: %q, dropped: %q)", op.ID, firstBodyParam, pr.Name))
					}
					firstBodyParam = pr.Name
				}

				pnames, ok := ptypes[pr.In]
				if !ok {
					pnames = make(map[string]struct{})
					ptypes[pr.In] = pnames
				}

				_, ok = pnames[pr.Name]
				if ok {
					res.AddErrors(errors.New(422, "duplicate parameter name %q for %q in operation %q", pr.Name, pr.In, op.ID))
				}
				pnames[pr.Name] = struct{}{}
			}
		}
	}
	return res
}

func parsePath(path string) ([]string, map[string]int) {
	return nil, nil
}

func (s *SpecValidator) validateReferencesValid() *Result {
	// each reference must point to a valid object
	return nil
}

func (s *SpecValidator) validateDefaultValueValidAgainstSchema() *Result {
	// every default value that is specified must validate against the schema for that property
	return nil
}