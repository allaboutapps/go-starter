package templates

type ViewTemplate string

const (
	ViewTemplateAccountConfirmation ViewTemplate = "account_confirmation.html.tmpl"
	ViewAppleSiteAssociation        ViewTemplate = "apple-app-site-association.json.tmpl"
	ViewAndroidAssetlinks           ViewTemplate = "assetlinks.json.tmpl"
)

func (vt ViewTemplate) String() string {
	return string(vt)
}
