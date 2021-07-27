module github.com/mories76/terraform-provider-ardoq

go 1.16

// require (
// 	github.com/mories76/ardoq-client-go v0.0.0 => ../ardoq-client-go
// 	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.0-rc.2
// )

require github.com/mories76/ardoq-client-go v0.0.0

replace github.com/mories76/ardoq-client-go v0.0.0 => ../ardoq-client-go

require (
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.0
)
