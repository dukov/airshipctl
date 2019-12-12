package render

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"sigs.k8s.io/kustomize/v3/pkg/fs"

	"opendev.org/airship/airshipctl/pkg/document"
)

// Render prints out filtered documents
func (s *Settings) Render(path string, out io.Writer) error {
	if err := s.verifyInput(); err != nil {
		return err
	}

	inputFilters := map[string][]string{
		labelJSONPathPrefix:      s.Label,
		annotationJSONPathPrefix: s.Annotation,
		kindJSONPathPrefix:       s.Kind,
		apiVersionJSONPathPrefix: s.GroupVersion,
	}

	docBundle, err := document.NewBundle(fs.MakeRealFS(), path, "")
	if err != nil {
		return err
	}
	filters := make([]string, 0)
	for prefix, args := range inputFilters {
		filters = append(filters, s.prepareFilter(prefix, args)...)
	}

	filterExpr := strings.Join(filters, " && ")
	if s.RawFilter != "" {
		filterExpr = s.RawFilter
	}

	if filterExpr == "" {
		return docBundle.Write(out)
	}

	filteredBundle, err := document.EvaluateExpressionFilter(filterExpr, docBundle)
	if err != nil {
		return err
	}

	return filteredBundle.Write(out)
}

func (s *Settings) prepareFilter(jPathPrefix string, filters []string) (res []string) {
	for _, filter := range filters {
		matchingOp := "=="
		if strings.HasPrefix(filter, "!") {
			matchingOp = "!="
			filter = strings.TrimPrefix(filter, "!")
		}

		filterField := jPathPrefix
		filterVal := filter

		filterKeyVal := strings.Split(filter, "=")
		if len(filterKeyVal) == 2 {
			filterField = fmt.Sprintf(`%s["%s"]`, jPathPrefix, filterKeyVal[0])
			filterVal = filterKeyVal[1]
		}

		res = append(res, filterField+matchingOp+strconv.Quote(filterVal))
	}
	return res
}

func (s *Settings) verifyInput() error {
	nonRaw := false
	for _, arg := range [...][]string{s.Label, s.Annotation, s.Kind, s.GroupVersion} {
		if arg != nil {
			nonRaw = true
			break
		}
	}
	if s.RawFilter != "" && nonRaw {
		return document.ErrWrongRenderArgs{}
	}

	for _, arg := range [...][]string{s.Label, s.Annotation} {
		for _, filter := range arg {
			if len(strings.Split(filter, "=")) != 2 {
				return document.ErrBadRenderArgFormat{Arg: filter}
			}
		}
	}
	return nil
}
