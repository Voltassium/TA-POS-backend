package validation

type Enum interface {
	IsValidEnum() bool
}

var BankChoices = []string{
	"Bank Mandiri",
	"Bank Negara Indonesia (BNI)",
	"Bank Central Asia (BCA)",
	"Bank Rakyat Indonesia (BRI)",
	"Bank Danamon",
	"CIMB Niaga",
	"Bank Permata",
	"OCBC NISP",
	"Bank Muamalat",
	"Bank Syariah Indonesia (BSI)",
	"Maybank Indonesia",
	"Bank Negara Indonesia Syariah (BNI Syariah)",
	"Bank Panin",
	"Bank Victoria International",
	"Bank Mega",
	"Bank Artha Graha",
	"Bank Sinarmas",
	"Bank Jabar Banten",
	"Bank BJB (Bank Jawa Barat dan Banten)",
	"Maybank",
}

const (
	CheckOnlyAlphabet = `^[A-Za-z\s._\-\/]+$`
	Digits            = `^\d+$`
	NPWP              = `^(\d{2}\.\d{3}\.\d{3}\.\d{1}\-\d{3})$`
	KTP               = `^\d{16}$`
)
