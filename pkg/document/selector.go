package document

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type Operation string

// SlectorField type representin JSON path to particular k8s resource
type SelectorField string

// constants for logical 'AND' amd logical 'OR' operations
const (
	LogicalAnd Operation = "&&"
	LogicalOr  Operation = "||"
)

const (
	Eq    Operation = "=="
	NotEq Operation = "!="
)

const (
	Kind       SelectorField = "kind"
	Label      SelectorField = "metadata.labels"
	Annotation SelectorField = "metadata.annotations"
	ApiVersion SelectorField = "apiVersion"

	metaFiledsCount = 4
)

type Selector interface {
	ToRaw() string
}

func NewKindSelector(val string, op selection.Operator) Selector {
	return &FieldSelector{
		FieldPath: string(Kind),
		Operation: op,
		Value:     val,
	}
}

func NewApiVersionSelector(val string, op selection.Operator) Selector {
	return &FieldSelector{
		FieldPath: string(ApiVersion),
		Operation: op,
		Value:     val,
	}
}

func NewLabelSelector(val string) (Selector, error) {
	return keyValSelector(Label, val)
}

func NewAnnotationSelector(val string) (Selector, error) {
	return keyValSelector(Annotation, val)
}

func keyValSelector(prefix SelectorField, val string) (Selector, error) {
	reqs, err := labels.ParseToRequirements(val)
	if err != nil {
		return nil, err
	}

	if len(reqs) != 1 {
		return nil, fmt.Errorf("Specify exactly ONE selector")
	}
	req := reqs[0]

	oper := req.Operator()
	switch oper {
	case selection.Equals, selection.DoubleEquals, selection.NotEquals:
		if oper == selection.Equals {
			oper = selection.DoubleEquals
		}
		if req.Values().Len() != 1 {
			return nil, fmt.Errorf("Specify exactly ONE value for selector")
		}
		return &FieldSelector{
			FieldPath: string(prefix) + "." + req.Key(),
			Operation: oper,
			Value:     req.Values().List()[0],
		}, nil
	}
	return nil, fmt.Errorf("Operator %s is not supported", oper)
}

type FieldSelector struct {
	FieldPath string
	Operation selection.Operator
	Value     interface{}
}

func (fs *FieldSelector) ToRaw() string {
	return fs.FieldPath + string(fs.Operation) + fmt.Sprintf("%v", fs.Value)
}

type Filter struct {
	Selectors []Selector
	Operator  Operation
}

func (f *Filter) ToRaw() string {
	rawSelectors := make([]string, 0, len(f.Selectors))
	for _, sel := range f.Selectors {
		rawSelectors = append(rawSelectors, sel.ToRaw())
	}
	return strings.Join(rawSelectors, string(f.Operator))
}
