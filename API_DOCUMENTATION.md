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
  - `cod_bod` (opcional): Filtrar por código de bodega
  - `page` (opcional): Número de página para paginación
  - `pageSize` (opcional): Cantidad de elementos por página (default: 100 cuando se usa paginación)

- **Comportamiento de paginación**:
  - Si no se proporcionan `page` o `pageSize`, se devuelven todos los resultados sin paginación
  - Si se proporciona al menos uno de estos parámetros, se activa la paginación

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

# Filtrado por código de artículo
GET http://localhost:8080/api/allsaldos?cod_art=IJUN5320BRO

# Filtrado por zeta artículo
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
  "Bod": "código de bodega",
  "Código_Artículo": "código del artículo",
  "Zeta_Articulo": "código zeta del artículo",
  "Descripción_Artículo": "descripción",
  "U_C": 0,
  "U_M": 0,
  "Cif": 0,
  "Costo": 0,
  "Precio_Vta": 0,
  "Saldo_Disponible": 0
}
```

## Notas

1. El endpoint filtra automáticamente para mostrar solo productos con saldo disponible positivo.
2. Las consultas de saldos por defecto muestran el año 2025.
3. La paginación es opcional. Cuando no se especifica, se devuelven todos los resultados en una sola respuesta.
4. El valor `Precio_Vta` corresponde al valor máximo disponible para cada producto.
5. El campo `Saldo_Disponible` se calcula como: `sal_ant + tot_ent - tot_sal - sal_com`
