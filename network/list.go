package network

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n", nw.Name, nw.IpRange, nw.Driver)

	}
	w.Flush()
}
