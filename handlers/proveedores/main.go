package proveedores

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/periface/checador/internals/appsheets"
	"github.com/periface/checador/internals/models"
	"github.com/periface/checador/internals/utils"
	"github.com/periface/checador/views"
)

var renderers = utils.NewRenderers()

type ProveedoresHandlers struct {
}

func NewProveedoresHandlers() *ProveedoresHandlers {
	return &ProveedoresHandlers{}
}
func (chh *ProveedoresHandlers) ProveedoresIndex(c echo.Context) error {
	rfcQuery := c.QueryParam("rfc")
	component := views.Index(rfcQuery)
	return renderers.RenderNoLayout(c, http.StatusOK, component)
}

func (chh *ProveedoresHandlers) BuscarProveedor(c echo.Context) error {
	rfcQuery := c.QueryParam("rfc")
	if rfcQuery == "" {
		slog.Error("error")
		component := views.Buscar(models.BuscarResponse{})
		return renderers.RenderNoLayout(c, http.StatusOK, component)
	} else {
		appsheetsInstance, err := appsheets.NewAppsheets()
		if err != nil {
			slog.Error("Instance appsheets")
			slog.Error(err.Error())
		}
		proveedorInfo := fetchProveedorInfo(rfcQuery, appsheetsInstance)
		fmt.Println(proveedorInfo)
		component := views.Buscar(proveedorInfo)

		return renderers.RenderNoLayout(c, http.StatusOK, component)
	}
}

func buscarProveedorEnAtcom(rfc string, instance *appsheets.Appsheets) ([]map[string]string, error) {
	query := `Filter(PADRON DE PROVEEDORES, [RFC]=${rfc})`
	query = strings.ReplaceAll(query, "${rfc}", rfc)
	return instance.Search("PADRON DE PROVEEDORES", models.AppSheetsPayload{
		Action: "Find",
		Properties: map[string]string{
			"Selector": query,
		},
	})
}

// iterate all maps and gets only the valid Props with their values
func getOnlyThisProps(inputList []map[string]string, validProps []string) []map[string]string {
	// Crear un set para búsqueda rápida de props válidas
	validSet := make(map[string]struct{})
	for _, prop := range validProps {
		validSet[prop] = struct{}{}
	}

	var result []map[string]string

	// Iterar sobre cada mapa de la lista
	for _, item := range inputList {
		filtered := make(map[string]string)
		for key, value := range item {
			if _, ok := validSet[key]; ok {
				filtered[key] = value
			}
		}
		result = append(result, filtered)
	}

	return result
}


func fetchProveedorInfo(rfcQuery string, appsheetsInstance *appsheets.Appsheets) models.BuscarResponse {
	datosDelProveedor, err := buscarProveedorEnAtcom(rfcQuery, appsheetsInstance)

	if err != nil {
		slog.Error(err.Error())
	}
	proveedorInfo := models.BuscarResponse{
		InformacionDelProveedor: getOnlyThisProps(datosDelProveedor, []string{
			"RFC",
			"RAZON SOCIAL",
			"NOMBRE DEL PROVEEDOR",
			"1ER. APELLIDO",
			"2O. APELLIDO",
			"GIRO",
			"FECHA ALTA",
			"FECHA VENCIMIENTO",
		}),
	}
	return proveedorInfo
}
