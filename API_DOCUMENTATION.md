# Documentación de API de Saldos

Esta documentación describe los endpoints disponibles en la API de Saldos, sus parámetros y ejemplos de uso.

## Índice

- [Estado del Servidor](#estado-del-servidor)
- [Obtener Todos los Saldos](#obtener-todos-los-saldos)
- [Estructura de Datos](#estructura-de-datos)

## Estado del Servidor

Verifica si la API está en funcionamiento.

- **URL**: `/health`
- **Método**: `GET`
- **Respuesta exitosa**:
  - Código: 200
  - Contenido: `{ "status": "OK" }`

### Ejemplo

```
GET http://localhost:8080/health
```

## Obtener Todos los Saldos

Recupera saldos con varios filtros opcionales y paginación opcional.

- **URL**: `/api/allsaldos`
- **Método**: `GET`
- **Parámetros de consulta**:
  - `anio` (opcional): Año del saldo (default: 2025)
  - `cod_art` (opcional): Filtrar por código de artículo
  - `codigo` (opcional): Filtrar por zeta artículo (código zeta)
  - `cod_bod` (opcional): Filtrar por código de bodega (default: "01")
  - `page` (opcional): Número de página para paginación
  - `pageSize` (opcional): Cantidad de elementos por página (default: 100 cuando se usa paginación)

- **Comportamiento de paginación**:
  - Si no se proporcionan `page` o `pageSize`, se devuelven todos los resultados sin paginación
  - Si se proporciona al menos uno de estos parámetros, se activa la paginación

- **Comportamiento de valores ponderados**:
  - Cuando se filtra por `cod_art` o `codigo`, se incluyen valores ponderados adicionales
  - Estos valores son `Cif_Prom_Ponderado` y `Precio_Vta_Ponderado`
  - Se calculan usando la cantidad ingresada (`can_ing`) como factor de ponderación

- **Respuesta exitosa (sin paginación)**:
  - Código: 200
  - Contenido: 
    ```json
    {
      "data": [array de objetos Saldo],
      "total": número total de registros devueltos,
      "filtros": {
        "anio": "2025",
        "cod_art": valor del filtro,
        "codigo": valor del filtro,
        "cod_bod": valor del filtro
      }
    }
    ```

- **Respuesta exitosa (con paginación)**:
  - Código: 200
  - Contenido: 
    ```json
    {
      "data": [array de objetos Saldo],
      "total": número total de registros que cumplen los criterios,
      "page": página actual,
      "pages": número total de páginas,
      "filtros": {
        "anio": "2025",
        "cod_art": valor del filtro,
        "codigo": valor del filtro,
        "cod_bod": valor del filtro
      }
    }
    ```

- **Respuesta cuando no hay resultados**:
  - Código: 200
  - Contenido: 
    ```json
    {
      "message": "No se encontraron saldos disponibles con los filtros proporcionados",
      "data": [],
      "total": 0
    }
    ```

### Ejemplos

```
# Todos los saldos del año 2025 (sin paginación)
GET http://localhost:8080/api/allsaldos

# Con paginación
GET http://localhost:8080/api/allsaldos?page=1&pageSize=20

# Filtrado por código de artículo (incluye valores ponderados)
GET http://localhost:8080/api/allsaldos?cod_art=IJUN5320BRO

# Filtrado por zeta artículo (incluye valores ponderados)
GET http://localhost:8080/api/allsaldos?codigo=101-22-007302-039

# Filtrado por código de bodega
GET http://localhost:8080/api/allsaldos?cod_bod=01

# Combinando filtros y paginación
GET http://localhost:8080/api/allsaldos?anio=2024&cod_bod=01&page=2&pageSize=50
```

## Estructura de Datos

### Objeto Saldo
```json
{
  "Año": 2025,
  "Bod": "01",
  "Código_Artículo": "código del artículo",
  "Zeta_Articulo": "código zeta del artículo",
  "Descripción_Artículo": "descripción",
  "U_C": 0,
  "U_M": 0,
  "Cantidad_Ingresada": 100,
  "Cif": 3.15,
  "Costo": 3.18,
  "Precio_Vta": 4.5,
  "Saldo_Disponible": 85,
  "Cif_Prom_Ponderado": 3.17,
  "Precio_Vta_Ponderado": 4.52
}
```

## Notas

1. El endpoint filtra automáticamente para mostrar solo productos con saldo disponible positivo.
2. Las consultas de saldos por defecto muestran el año 2025.
3. La bodega predeterminada es la "01".
4. La paginación es opcional. Cuando no se especifica, se devuelven todos los resultados en una sola respuesta.
5. El valor `Precio_Vta` corresponde al valor máximo disponible para cada producto.
6. Los valores ponderados (`Cif_Prom_Ponderado` y `Precio_Vta_Ponderado`) solo aparecen cuando se filtra por un producto específico.
7. El campo `Saldo_Disponible` se calcula como: `sal_ant + tot_ent - tot_sal - sal_com`
