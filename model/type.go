package domyApi

import (
	"time"
)

type Token struct {
	Key    string
	Values string
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
	UserID    string    `bson:"user_id" json:"user_id"`
	Password  string    `bson:"password" json:"password"`
	Token     string    `bson:"token" json:"token"`
	Role      string    `bson:"role" json:"role"`
	NoHp      string    `bson:"nohp" json:"nohp"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Mahasiswa struct {
	Email            string `bson:"email,omitempty" json:"email,omitempty"`
	NIM              string `bson:"nim,omitempty" json:"nim,omitempty"`
	Nama             string `bson:"nama,omitempty" json:"nama,omitempty"`
	ProgramStudi     string `bson:"program_studi,omitempty" json:"program_studi,omitempty"`
	NomorHp          string `bson:"no_hp,omitempty" json:"no_hp,omitempty"`
	NIRM             string `bson:"nirm,omitempty" json:"nirm,omitempty"`
	PeriodeMasuk     string `bson:"periode_masuk,omitempty" json:"periode_masuk,omitempty"`
	TahunKurikulum   string `bson:"tahun_kurikulum,omitempty" json:"tahun_kurikulum,omitempty"`
	SistemKuliah     string `bson:"sistem_kuliah,omitempty" json:"sistem_kuliah,omitempty"`
	Kelas            string `bson:"kelas,omitempty" json:"kelas,omitempty"`
	JenisPendaftaran string `bson:"jenis_pendaftaran,omitempty" json:"jenis_pendaftaran,omitempty"`
	JalurPendaftaran string `bson:"jalur_pendaftaran,omitempty" json:"jalur_pendaftaran,omitempty"`
	Gelombang        string `bson:"gelombang,omitempty" json:"gelombang,omitempty"`
	TanggalMasuk     string `bson:"tanggal_masuk,omitempty" json:"tanggal_masuk,omitempty"`
	KebutuhanKhusus  string `bson:"kebutuhan_khusus,omitempty" json:"kebutuhan_khusus,omitempty"`
	StatusMahasiswa  string `bson:"status_mahasiswa,omitempty" json:"status_mahasiswa,omitempty"`
}

type Dosen struct {
	Email         string `bson:"email,omitempty" json:"email,omitempty"`
	NIP           string `bson:"nip,omitempty" json:"nip,omitempty"`
	NIDN          string `bson:"nidn,omitempty" json:"nidn,omitempty"`
	Nama          string `bson:"nama,omitempty" json:"nama,omitempty"`
	GelarDepan    string `bson:"gelar_depan,omitempty" json:"gelar_depan,omitempty"`
	GelarBelakang string `bson:"gelar_belakang,omitempty" json:"gelar_belakang,omitempty"`
	JenisKelamin  string `bson:"jenis_kelamin,omitempty" json:"jenis_kelamin,omitempty"`
	TempatLahir   string `bson:"tempat_lahir,omitempty" json:"tempat_lahir,omitempty"`
	TanggalLahir  string `bson:"tanggal_lahir,omitempty" json:"tanggal_lahir,omitempty"`
	Agama         string `bson:"agama,omitempty" json:"agama,omitempty"`
	NoHp          string `bson:"no_hp,omitempty" json:"no_hp,omitempty"`
	EmailKampus   string `json:"email_kampus,omitempty"`
	EmailPribadi  string `json:"email_pribadi,omitempty"`
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

type ListBimbingan struct {
	No              string `bson:"no,omitempty" json:"no,omitempty"`
	Tanggal         string `bson:"Tanggal,omitempty" json:"Tanggal,omitempty"`
	DosenPembimbing string `bson:"dosenpembimbing,omitempty" json:"dosenpembimbing,omitempty"`
	Topik           string `bson:"topik,omitempty" json:"topik,omitempty"`
	Disetujui       bool   `bson:"disetujui,omitempty" json:"disetujui,omitempty"`
	DataID          string `bson:"data_id,omitempty" json:"data_id,omitempty"`
}

type TugasAkhirAllMahasiswa struct {
	Nama         string `bson:"nama,omitempty" json:"nama"`
	NIM          string `bson:"nim,omitempty" json:"nim"`
	Judul        string `bson:"judul,omitempty" json:"judul"`
	Pembimbing1  string `bson:"pembimbing1,omitempty" json:"pembimbing1"`
	Pembimbing2  string `bson:"pembimbing2,omitempty" json:"pembimbing2"`
	TanggalMulai string `bson:"tanggal_mulai,omitempty" json:"tanggal_mulai"`
	Status       string `bson:"status,omitempty" json:"status"`
	DataID       string `bson:"data_id,omitempty" json:"data_id"`
}

type TugasAkhirMahasiswa struct {
	DataID       string `bson:"data_id,omitempty" json:"data_id"`
	Judul        string `bson:"judul,omitempty" json:"judul"`
	Pembimbing1  string `bson:"pembimbing1,omitempty" json:"pembimbing1"`
	Pembimbing2  string `bson:"pembimbing2,omitempty" json:"pembimbing2"`
	TanggalMulai string `bson:"tanggal_mulai,omitempty" json:"tanggal_mulai"`
	Status       string `bson:"status,omitempty" json:"status"`
}

type JadwalMengajar struct {
	No           string `json:"no"`
	Kode         string `json:"kode"`
	MataKuliah   string `json:"mata_kuliah"`
	SKS          string `json:"sks"`
	Smt          string `json:"smt"`
	Kelas        string `json:"kelas"`
	ProgramStudi string `json:"program_studi"`
	Hari         string `json:"hari"`
	Waktu        string `json:"waktu"`
	Ruang        string `json:"ruang"`
}
