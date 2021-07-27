// how to build and copy
go build -o terraform-provider-ardoq
mv terraform-provider-ardoq ~/.terraform.d/plugins/mories.com/terraform/ardoq/0.1/darwin_arm64

// test terraform
terraform -chdir=./examples/components validate

// omgevingsvariablen die gezet moeten zijn voor de ardoq provider
export ARDOQ_APIKEY=
export ARDOQ_BASEURI=
export ARDOQ_ORG=



// voorbereiding op macOS om te zorgen dat de debugger in vscode mag verbinden met bestaande processen
sudo /usr/sbin/DevToolsSecurity --enable
sudo dscl . append /Groups/_developer GroupMembership mories

// provider in debug modus starten
// https://www.terraform.io/docs/extend/debugging.html#starting-a-provider-in-debug-mode

/Users/mories/go/bin/dlv exec --listen=127.0.0.1:58482 --headless ~/.terraform.d/plugins/mories.com/terraform/ardoq/0.1/darwin_arm64/terraform-provider-ardoq -- --debug 