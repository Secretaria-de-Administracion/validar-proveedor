package handlers

import (
	"github.com/periface/checador/handlers/proveedores"
)

type ChecadorHttpHanlder struct {
	Proveedores *proveedores.ProveedoresHandlers
}

func NewMainHandler() *ChecadorHttpHanlder {
	proveedores := proveedores.NewProveedoresHandlers()

	return &ChecadorHttpHanlder{
		Proveedores: proveedores,
	}
}
