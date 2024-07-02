package domyApi

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type Bimbingan struct {
	Bimbinganke    string `bson:"bimbinganke,omitempty" json:"bimbinganke,omitempty"`
	NIP            string `bson:"nip,omitempty" json:"nip,omitempty"`
	TglBimbingan   string `bson:"tglbimbingan,omitempty" json:"tglbimbingan,omitempty"`
	TopikBimbingan string `bson:"topikbimbingan,omitempty" json:"topikbimbingan,omitempty"`
	Bahasan        string `bson:"bahasan,omitempty" json:"bahasan,omitempty"`
	Link           string `bson:"link,omitempty" json:"link,omitempty"`
	Lampiran       string `bson:"lampiran,omitempty" json:"lampiran,omitempty"`
	Key            string `bson:"key,omitempty" json:"key,omitempty"`
	Act            string `bson:"act,omitempty" json:"act,omitempty"`
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

type ResponseAct struct {
	Login     bool   `json:"login"`
	SxSession string `json:"token"`
}

type RequestLoginSiakad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ResponseLogin struct {
	Code    string `json:"code"`
	Session string `json:"session"`
	Role    string `json:"role"`
}

type TokenData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Token     string             `bson:"token" json:"token"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
