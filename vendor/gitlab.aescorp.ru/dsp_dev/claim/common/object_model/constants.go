package object_model

const (
	// DocTypeInvoice -- тип платежного документа С/Ф
	DocTypeInvoice = 35
)

type crud_transport_type int

const (
	NeedUseNRPC crud_transport_type = iota
	NeedUseDB
	NeedUseGRPC
)

var Crud_transport crud_transport_type = 0
