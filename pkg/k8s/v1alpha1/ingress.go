package v1alpha1

type Ingress struct {
	Namespace   string
	Host        string
	Path        string
	PathType    string
	ServiceName string
	ServicePort int
}
