package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pablojnd/api_go_saldo/database"
	"github.com/pablojnd/api_go_saldo/models"
)

// GetSaldos obtiene los saldos con filtros opcionales
func GetSaldos(c *gin.Context) {
	var saldos []models.Saldo

	query := database.DB.Model(&models.Saldo{})

	// Aplicar filtros si se proporcionan
	codArt := c.Query("codArt")
	if codArt != "" {
		query = query.Where("cod_art = ?", codArt)
	} else {
		// Si no se especifica un código, usar el de ejemplo
		query = query.Where("cod_art = ?", "IJUN5320BRO")
	}

	anio := c.DefaultQuery("anio", "2025")
	query = query.Where("anio_pro = ?", anio)

	// Filtro opcional por código de bodega
	if codBod := c.Query("codBod"); codBod != "" {
		query = query.Where("cod_bod = ?", codBod)
	}

	// Filtro opcional por saldo mínimo
	if minSaldo := c.Query("minSaldo"); minSaldo != "" {
		minVal, err := strconv.ParseFloat(minSaldo, 64)
		if err == nil {
			query = query.Where("(sal_ant + tot_ent - tot_sal - sal_com) >= ?", minVal)
		}
	} else {
		// Aplicar la condición para saldo disponible
		query = query.Where("(sal_ant + tot_ent - tot_sal - sal_com) > 0")
	}

	// Paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// Contar total de registros para paginación
	var total int64
	query.Count(&total)

	// Seleccionar campos y aplicar cálculo para saldo disponible
	query = query.Select("anio_pro, cod_bod, cod_art, zet_art, des_adu, uni_set, uni_caj, " +
		"cif_uni, cos_rea, val_viu, (sal_ant + tot_ent - tot_sal - sal_com) as Saldo_Disponible, " +
		"sal_ant, tot_ent, tot_sal, sal_com")

	// Aplicar paginación
	if pageSize > 0 {
		query = query.Limit(pageSize).Offset(offset)
	}

	// Agrupar y ordenar
	query = query.Group("zet_art, anio_pro").Order("zet_art")

	result := query.Find(&saldos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener saldos",
			"details": result.Error.Error(),
		})
		return
	}

	// Si no se encontraron registros
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No se encontraron saldos que cumplan con los criterios",
			"data":    []models.Saldo{},
			"total":   0,
			"page":    page,
			"pages":   0,
		})
		return
	}

	// Calcular el campo saldo disponible para cada resultado
	for i := range saldos {
		saldos[i].SaldoDisponible = saldos[i].SalAnt + saldos[i].TotEnt - saldos[i].TotSal - saldos[i].SalCom
	}

	// Calcular páginas totales
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// Respuesta paginada
	c.JSON(http.StatusOK, gin.H{
		"data":  saldos,
		"total": total,
		"page":  page,
		"pages": totalPages,
	})
}

// GetSaldoById obtiene un saldo específico por código de artículo
func GetSaldoById(c *gin.Context) {
	codArt := c.Param("codArt")
	if codArt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de artículo no proporcionado"})
		return
	}

	var saldos []models.Saldo
	query := database.DB.Model(&models.Saldo{}).
		Where("cod_art = ?", codArt).
		Where("(sal_ant + tot_ent - tot_sal - sal_com) > 0").
		Select("anio_pro, cod_bod, cod_art, zet_art, des_adu, uni_set, uni_caj, " +
			"cif_uni, cos_rea, val_viu, (sal_ant + tot_ent - tot_sal - sal_com) as Saldo_Disponible, " +
			"sal_ant, tot_ent, tot_sal, sal_com").
		Group("zet_art, anio_pro").
		Order("zet_art")

	if result := query.Find(&saldos); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el saldo"})
		return
	}

	if len(saldos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró el artículo o no tiene saldo disponible"})
		return
	}

	// Calcular el campo saldo disponible para cada resultado
	for i := range saldos {
		saldos[i].SaldoDisponible = saldos[i].SalAnt + saldos[i].TotEnt - saldos[i].TotSal - saldos[i].SalCom
	}

	c.JSON(http.StatusOK, saldos)
}

// GetSaldos2025 obtiene los saldos disponibles del año 2025 con paginación
func GetSaldos2025(c *gin.Context) {
	type SaldoResponse struct {
		Anio            int     `json:"Año"`
		Bod             string  `json:"Bod"`
		CodigoArticulo  string  `json:"Código_Artículo"`
		ZetaArticulo    string  `json:"Zeta_Articulo"`
		DescArticulo    string  `json:"Descripción_Artículo"`
		UC              int     `json:"U_C"`
		UM              float64 `json:"U_M"`
		Cif             float64 `json:"Cif"`
		Costo           float64 `json:"Costo"`
		PrecioVta       float64 `json:"Precio_Vta"`
		SaldoDisponible float64 `json:"Saldo_Disponible"`
	}

	// Paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	offset := (page - 1) * pageSize

	// Primero obtener el conteo total para la paginación
	var total int64
	countQuery := `
		SELECT COUNT(*) 
		FROM (
			SELECT 1
			FROM saldos
			WHERE anio_pro = 2025
			AND (sal_ant + tot_ent - tot_sal - sal_com) > 0
			GROUP BY zet_art, anio_pro
		) as count_table
	`
	database.DB.Raw(countQuery).Count(&total)

	// Consulta optimizada con paginación
	var saldos []SaldoResponse
	query := `
		SELECT 
			anio_pro AS Anio,
			cod_bod AS Bod,
			cod_art AS CodigoArticulo,
			zet_art AS ZetaArticulo,
			des_adu AS DescArticulo,
			uni_set AS UC,
			uni_caj AS UM,
			cif_uni AS Cif,
			cos_rea AS Costo,
			val_viu AS PrecioVta,
			(sal_ant + tot_ent - tot_sal - sal_com) AS SaldoDisponible
		FROM 
			saldos
		WHERE 
			anio_pro = 2025
			AND (sal_ant + tot_ent - tot_sal - sal_com) > 0
		GROUP BY 
			zet_art, anio_pro
		ORDER BY 
			zet_art
		LIMIT ? OFFSET ?
	`

	if err := database.DB.Raw(query, pageSize, offset).Scan(&saldos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener saldos 2025",
			"details": err.Error(),
		})
		return
	}

	// Si no se encontraron registros
	if len(saldos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No se encontraron saldos disponibles para 2025",
			"data":    []SaldoResponse{},
			"total":   0,
			"page":    page,
			"pages":   0,
		})
		return
	}

	// Calcular páginas totales
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  saldos,
		"total": total,
		"page":  page,
		"pages": totalPages,
	})
}

// GetSumaSaldoProducto obtiene un producto con la suma del Saldo_Disponible
func GetSumaSaldoProducto(c *gin.Context) {
	codArt := c.Param("codArt")
	if codArt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código de artículo no proporcionado"})
		return
	}

	type SumaSaldo struct {
		CodArt         string  `json:"codigo_articulo"`
		ZetArt         string  `json:"zeta_articulo"`
		DesAdu         string  `json:"descripcion"`
		SumaSaldo      float64 `json:"suma_saldo_disponible"`
		PrecioVenta    float64 `json:"precio_venta"`
		CostoPonderado float64 `json:"costo_ponderado"`
	}

	var resultado SumaSaldo

	// Usar SQL directo para calcular el costo ponderado correctamente
	query := `
		SELECT 
			cod_art AS CodArt, 
			zet_art AS ZetArt, 
			des_adu AS DesAdu, 
			SUM(sal_ant + tot_ent - tot_sal - sal_com) AS SumaSaldo, 
			MAX(val_viu) AS PrecioVenta, 
			SUM(cos_rea * (sal_ant + tot_ent - tot_sal - sal_com)) / SUM(sal_ant + tot_ent - tot_sal - sal_com) AS CostoPonderado
		FROM 
			saldos
		WHERE 
			cod_art = ? 
			AND (sal_ant + tot_ent - tot_sal - sal_com) > 0
		GROUP BY 
			cod_art, zet_art, des_adu
	`

	if err := database.DB.Raw(query, codArt).Scan(&resultado).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la suma del saldo"})
		return
	}

	if resultado.CodArt == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró el artículo o no tiene saldo disponible"})
		return
	}

	c.JSON(http.StatusOK, resultado)
}
