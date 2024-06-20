package domyApi

type Token struct {
	Key    string
	Values string
}

type Mahasiswa struct {
	NIM          string `bson:"nim"`
	Nama         string `bson:"nama"`
	ProgramStudi string `bson:"program_studi"`
	NoHp         string `bson:"no_hp"`
}

type Dosen struct {
	NIP  string `bson:"nip"`
	NIDN string `bson:"nidn"`
	Nama string `bson:"nama"`
	NoHp string `bson:"no_hp"`
}

type Profile struct {
	Token       string `bson:"token"`
	Phonenumber string `bson:"phonenumber"`
	Secret      string `bson:"secret"`
	URL         string `bson:"url"`
	QRKeyword   string `bson:"qrkeyword"`
	PublicKey   string `bson:"publickey"`
}

type Response struct {
	Response string `json:"response"`
	Info     string `json:"info,omitempty"`
	Status   string `json:"status,omitempty"`
	Location string `json:"location,omitempty"`
}
