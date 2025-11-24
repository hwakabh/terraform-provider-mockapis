terraform {
  required_providers {
    mockapis = {
      source = "registry.terraform.io/hwakabh/mockapis"
    }
  }
}

provider "mockapis" {
  username = "hwakabh"
  token    = "dummy"
}

data "mockapis_me" "this" {}

output "me_responses" {
  value = data.mockapis_me.this
}
