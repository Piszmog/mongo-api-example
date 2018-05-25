package model

type CFEnv struct {
    Mlab []MLab `json:"mlab"`
}
type MLab struct {
    Name             string      `json:"name"`
    Instance_name    string      `json:"instance_name"`
    Binding_name     string      `json:"binding_name"`
    Credentials      Credentials `json:"credentials"`
    Syslog_drain_url string      `json:"syslog_drain_url"`
    Volume_mounts    []string    `json:"volume_mounts"`
    Label            string      `json:"label"`
    Provider         string      `json:"provider"`
    Plan             string      `json:"plan"`
    Tags             []string    `json:"tags"`

}
type Credentials struct {
    Uri string `json:"uri"`
}
