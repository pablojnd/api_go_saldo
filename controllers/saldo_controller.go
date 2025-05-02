package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myuser/api/database"
	"github.com/myuser/api/models"
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
