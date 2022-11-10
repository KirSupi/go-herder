package herder

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"text/tabwriter"
)

func StringifyStates(states []ProcessState) string {
	needToPrintLabel := false
	for _, s := range states {
		if s.Label != nil {
			needToPrintLabel = true
			break
		}
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', tabwriter.Debug)
	if needToPrintLabel {
		_, _ = fmt.Fprint(w, "ID\tLabel\tActive")
	} else {
		_, _ = fmt.Fprint(w, "ID\tActive")
	}
	for _, s := range states {
		if needToPrintLabel {
			if s.Label != nil {
				_, _ = fmt.Fprintf(w, "\n%d\t%s\t%v", s.ID, *s.Label, s.Active)
			} else {
				_, _ = fmt.Fprintf(w, "\n%d\t\t%v", s.ID, s.Active)
			}
		} else {
			_, _ = fmt.Fprintf(w, "\n%d\t%v", s.ID, s.Active)
		}
	}
	_ = w.Flush()
	return b.String()
}

func errorNoProcessID(id int) error {
	return errors.New("no process with id #" + strconv.Itoa(id))
}
