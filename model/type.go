package domyApi

type Token struct {
	Key    string
	Values string
}

type Mahasiswa struct {
	NIM          string `bson:"nim,omitempty" json:"nim,omitempty"`
	Nama         string `bson:"nama,omitempty" json:"nama,omitempty"`
	ProgramStudi string `bson:"program_studi,omitempty" json:"program_studi,omitempty"`
	NomorHp      string `bson:"no_hp,omitempty" json:"no_hp,omitempty"`
}

type Dosen struct {
	NIP  string `bson:"nip,omitempty" json:"nip,omitempty"`
	NIDN string `bson:"nidn,omitempty" json:"nidn,omitempty"`
	Nama string `bson:"nama,omitempty" json:"nama,omitempty"`
	NoHp string `bson:"no_hp,omitempty" json:"no_hp,omitempty"`
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
