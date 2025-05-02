package models

// Saldo representa la estructura de la tabla saldos
type Saldo struct {
	AnioPro         int     `json:"Año" gorm:"column:anio_pro"`
	CodBod          string  `json:"Bod" gorm:"column:cod_bod"`
	CodArt          string  `json:"Código_Artículo" gorm:"column:cod_art"`
	ZetArt          string  `json:"Zeta_Articulo" gorm:"column:zet_art"`
	DesAdu          string  `json:"Descripción_Artículo" gorm:"column:des_adu"`
	UniSet          int     `json:"U_C" gorm:"column:uni_set"`
	UniCaj          float64 `json:"U_M" gorm:"column:uni_caj"`
	CifUni          float64 `json:"Cif" gorm:"column:cif_uni"`
	CosRea          float64 `json:"Costo" gorm:"column:cos_rea"`
	ValViu          float64 `json:"Precio_Vta" gorm:"column:val_viu"`
	SaldoDisponible float64 `json:"Saldo_Disponible" gorm:"-"` // Campo calculado
	SalAnt          float64 `json:"-" gorm:"column:sal_ant"`   // Usado para cálculo
	TotEnt          float64 `json:"-" gorm:"column:tot_ent"`   // Usado para cálculo
	TotSal          float64 `json:"-" gorm:"column:tot_sal"`   // Usado para cálculo
	SalCom          float64 `json:"-" gorm:"column:sal_com"`   // Usado para cálculo
}

// TableName especifica el nombre de la tabla
func (Saldo) TableName() string {
	return "saldos"
}
