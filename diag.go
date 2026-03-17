package ptrstruct

// FormatDiagnostic produces the diagnostic message for a violation.
// position is a human label such as "receiver", "parameter req", "field Meta".
func FormatDiagnostic(position string, v *Violation) string {
	if v.Path == "" {
		return position + " uses value struct " + v.TypeName + "; use *" + v.TypeName
	}
	return position + " uses " + v.Path + " " + v.TypeName + " by value"
}
