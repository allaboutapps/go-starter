package templates

type ViewTemplate string

const (
	ViewTemplateAccountConfirmation ViewTemplate = "account_confirmation.html.tmpl"
)

func (vt ViewTemplate) String() string {
	return string(vt)
}
