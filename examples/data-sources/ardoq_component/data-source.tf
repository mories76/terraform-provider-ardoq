data "ardoq-component" "someserver" {
    root_workspace = "<id>"
    name = "SomeServerName"
}

output "someserver" {
    value = data.ardoq_component.someserver
}