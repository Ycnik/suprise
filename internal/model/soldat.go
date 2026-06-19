package model

import "time"

type Soldat struct {
	ID           int          `json:"id" gorm:"primaryKey;column:id"`
	Version      int          `json:"version" gorm:"column:version"`
	Vorname      string       `json:"vorname" gorm:"column:vorname"`
	Nachname     string       `json:"nachname" gorm:"column:nachname"`
	Geburtsdatum *time.Time   `json:"geburtsdatum,omitempty" gorm:"column:geburtsdatum;type:date"`
	Geschlecht   *string      `json:"geschlecht,omitempty" gorm:"column:geschlecht;type:geschlecht"`
	Rang         *string      `json:"rang,omitempty" gorm:"column:rang;type:rang"`
	Username     string       `json:"username" gorm:"column:username"`
	Erzeugt      time.Time    `json:"erzeugt" gorm:"column:erzeugt"`
	Aktualisiert time.Time    `json:"aktualisiert" gorm:"column:aktualisiert"`
	Ausruestung  *Ausruestung `json:"ausruestung,omitempty" gorm:"foreignKey:SoldatID"`
	Verletzungen []Verletzung `json:"verletzungen,omitempty" gorm:"foreignKey:SoldatID"`
}

func (Soldat) TableName() string {
	return "soldat.soldat"
}

type Ausruestung struct {
	ID           int    `json:"id" gorm:"primaryKey;column:id"`
	Waffe        string `json:"waffe" gorm:"column:waffe;type:waffe"`
	Seriennummer string `json:"seriennummer" gorm:"column:seriennummer"`
	SoldatID     int    `json:"soldat_id" gorm:"column:soldat_id"`
}

func (Ausruestung) TableName() string {
	return "soldat.ausruestung"
}

type Verletzung struct {
	ID                     int       `json:"id" gorm:"primaryKey;column:id"`
	Verletzungsbezeichnung string    `json:"verletzungsbezeichnung" gorm:"column:verletzungsbezeichnung"`
	Behandelt              bool      `json:"behandelt" gorm:"column:behandelt"`
	Schweregrad            string    `json:"schweregrad" gorm:"column:schweregrad;type:schweregrad"`
	Verletzungsdatum       time.Time `json:"verletzungsdatum" gorm:"column:verletzungsdatum;type:date"`
	SoldatID               int       `json:"soldat_id" gorm:"column:soldat_id"`
}

func (Verletzung) TableName() string {
	return "soldat.verletzung"
}
