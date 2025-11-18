module 5tk.dev/virt

go 1.25.3

require (
	5tk.dev/ip v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	gopkg.in/yaml.v3 v3.0.1
)

require 5tk.dev/c3po v0.1.0 // indirect

replace 5tk.dev/c3po => ../c3po

replace 5tk.dev/ip => ../ip
