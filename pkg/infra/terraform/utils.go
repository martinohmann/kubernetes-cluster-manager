package terraform

const terraformCmd = "terraform"

func terraform(args ...string) []string {
	return append([]string{terraformCmd}, args...)
}
