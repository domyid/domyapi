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
