package user

import "time"

type User struct {
	ID           int		`gorm:"primaryKey;autoIncrement"`
	Name         string		`gorm:"size:255"`
	PasswordHash string		`gorm:"size:255"`
	Email        string		`gorm:"size:255"`
	Phone		string		`gorm:"size:50"`
	Avatar       string		`gorm:"size:255"`
	Role         string		`gorm:"size:10"`
	CreatedAt    time.Time	`gorm:"autoCreateTime"`
	UpdatedAt    time.Time	`gorm:"autoUpdateTime"`
	Addres		 Addres		`gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}


type Addres struct{
	ID 			int		`gorm:"primaryKey;autoIncrement"`
	UserID		int		`gorm:"not null;index"`
	Provinsi 	string	`gorm:"size:100"`
	Kabupaten 	string	`gorm:"size:100"`
	Kecamatan 	string	`gorm:"size:100"`
	Desa 		string	`gorm:"size:100"`
	Jalan 		string	`gorm:"size:100"`
	Rt 			string	`gorm:"size:100"`
	Rw 			string	`gorm:"size:100"`
	NoRm 		string	`gorm:"size:100"`
	KodePos 	string	`gorm:"size:100"`
	Deskripsi 	string	`gorm:"type:text"`
}