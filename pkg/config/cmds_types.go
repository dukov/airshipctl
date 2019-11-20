package config

type ClusterOptions struct {
	Name                  string
	ClusterType           string
	Server                string
	InsecureSkipTLSVerify bool
	CertificateAuthority  string
	EmbedCAData           bool
}

type ContextOptions struct {
	Name           string
	CurrentContext bool
	Cluster        string
	ClusterType    string
	AuthInfo       string
	Manifest       string
	Namespace      string
}
