package unindent

import "strings"

func Unindent(s string) string {
        lines := strings.Split(s, "\n")
        output := strings.Builder{}
        for _, l := range lines {
                output.WriteString(strings.TrimLeft(l, " \t"))
        }
        return output.String()
}
