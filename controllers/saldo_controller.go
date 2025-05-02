package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pablojnd/api_go_saldo/database"
)

// GetAllSaldos obtiene todos los saldos disponibles con opciones de filtrado y paginación
func GetAllSaldos(c *gin.Context) {
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

	// Parámetros de filtro
	anio := c.DefaultQuery("anio", "2025")
	codArt := c.Query("cod_art")
	zetaArt := c.Query("codigo")
	codBod := c.Query("cod_bod")

	// Parámetros de paginación
	usePagination := c.Query("page") != "" || c.Query("pageSize") != ""
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	offset := (page - 1) * pageSize

	// Construir condiciones WHERE
	whereConditions := "anio_pro = ? AND (sal_ant + tot_ent - tot_sal - sal_com) > 0"
	args := []interface{}{anio}

	// Agregar filtros adicionales si se proporcionaron
	if codArt != "" {
		whereConditions += " AND cod_art = ?"
		args = append(args, codArt)
	}

	if zetaArt != "" {
		whereConditions += " AND zet_art = ?"
		args = append(args, zetaArt)
	}

	if codBod != "" {
		whereConditions += " AND cod_bod = ?"
		args = append(args, codBod)
	}

	// Contar total de registros para paginación
	var total int64
	if usePagination {
		countQuery := `
			SELECT COUNT(*) 
			FROM (
				SELECT 1
				FROM saldos
				WHERE ` + whereConditions + `
				GROUP BY zet_art, anio_pro
			) as count_table
		`
		database.DB.Raw(countQuery, args...).Count(&total)
	}

	// Construir consulta principal
	mainQuery := `
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
			MAX(val_viu) AS PrecioVta,
			(sal_ant + tot_ent - tot_sal - sal_com) AS SaldoDisponible
		FROM 
			saldos
		WHERE 
			` + whereConditions + `
		GROUP BY 
			zet_art, anio_pro
		ORDER BY 
			zet_art
	`

	// Agregar paginación si se solicitó
	if usePagination {
		mainQuery += " LIMIT ? OFFSET ?"
		args = append(args, pageSize, offset)
	}

	// Ejecutar consulta
	var saldos []SaldoResponse
	if err := database.DB.Raw(mainQuery, args...).Scan(&saldos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al obtener los saldos",
			"details": err.Error(),
		})
		return
	}

	// Si no se encontraron registros
	if len(saldos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No se encontraron saldos disponibles con los filtros proporcionados",
			"data":    []SaldoResponse{},
			"total":   0,
		})
		return
	}

	// Preparar respuesta
	response := gin.H{
		"data":  saldos,
		"total": len(saldos),
		"filtros": gin.H{
			"anio":    anio,
			"cod_art": codArt,
			"codigo":  zetaArt,
			"cod_bod": codBod,
		},
	}

	// Agregar información de paginación si se está usando
	if usePagination {
		totalPages := int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}

		response["page"] = page
		response["pages"] = totalPages
		response["total"] = total
	}

	c.JSON(http.StatusOK, response)
}
